package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Environment variables controlling the deliberately-staged memory leak.
//
// The leak is OFF unless memoryLeakEnabledEnv is set to a truthy value
// (as understood by strconv.ParseBool), so the default behaviour of the
// application is completely unaffected. This is intentionally an opt-in test
// fixture for reproducing an OOMKill-on-a-loop failure, not a real bug.
const (
	memoryLeakEnabledEnv = "MEMORY_LEAK_ENABLED"
	memoryLeakRateEnv    = "MEMORY_LEAK_RATE_BYTES_PER_SEC"
	memoryLeakGraceEnv   = "MEMORY_LEAK_GRACE_PERIOD"
)

const (
	// defaultMemoryLeakRateBytesPerSec is used when the leak is enabled but no
	// valid rate is configured. 1 MiB/s.
	defaultMemoryLeakRateBytesPerSec = 1 << 20
	// defaultMemoryLeakGracePeriod gives the pod time to pass its readiness
	// probe and reach Ready before memory starts climbing.
	defaultMemoryLeakGracePeriod = 30 * time.Second
	// leakPageSize is the stride, in bytes, at which freshly allocated memory
	// is written to. Touching one byte per OS page faults the whole page in so
	// it becomes resident; allocated-but-untouched memory would not show up in
	// container_memory_working_set_bytes.
	leakPageSize = 4096
	// leakInterval is how often memory is allocated. A sub-second interval
	// keeps the growth curve smooth and linear rather than steppy/bursty.
	leakInterval = 100 * time.Millisecond
)

// memoryLeaker holds the state of the running leak so it can be inspected and
// controlled at runtime.
type memoryLeaker struct {
	rateBytesPerSec int64
	gracePeriod     time.Duration

	// active reports whether allocation is currently happening. It can be
	// toggled at runtime via the control endpoints without a redeploy.
	active atomic.Bool

	mu     sync.Mutex
	blocks [][]byte // retained references are what keep the pages resident
	leaked int64

	gauge prometheus.Gauge
}

// setupMemoryLeak wires up the deliberately-staged memory leak when it is
// enabled via the environment.
//
// When disabled (the default) it is a no-op and leaves the application
// byte-for-byte identical to its previous behaviour: no goroutine is started,
// no metric is registered, no control endpoints are added and nothing is
// logged. It must be called before the HTTP server starts, as it registers
// handlers on http.DefaultServeMux.
func setupMemoryLeak() {
	enabled, _ := strconv.ParseBool(os.Getenv(memoryLeakEnabledEnv))
	if !enabled {
		return
	}

	rate := int64(defaultMemoryLeakRateBytesPerSec)
	if raw := os.Getenv(memoryLeakRateEnv); raw != "" {
		if v, err := strconv.ParseInt(raw, 10, 64); err == nil && v > 0 {
			rate = v
		} else {
			slog.Warn("Invalid memory leak rate, using default",
				"env", memoryLeakRateEnv, "value", raw,
				"default_bytes_per_sec", defaultMemoryLeakRateBytesPerSec)
		}
	}

	grace := time.Duration(defaultMemoryLeakGracePeriod)
	if raw := os.Getenv(memoryLeakGraceEnv); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil && d >= 0 {
			grace = d
		} else {
			slog.Warn("Invalid memory leak grace period, using default",
				"env", memoryLeakGraceEnv, "value", raw,
				"default", defaultMemoryLeakGracePeriod.String())
		}
	}

	l := &memoryLeaker{
		rateBytesPerSec: rate,
		gracePeriod:     grace,
		gauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "helloworld_memory_leaked_bytes",
			Help: "Bytes deliberately leaked by the staged memory-leak injector. Only present when the leak is enabled.",
		}),
	}

	// Make it obvious to anyone inspecting the environment that this failure is
	// staged, not a real bug.
	slog.Warn("Memory leak injection is DELIBERATELY ENABLED (staged failure for demos/testing)",
		"rate_bytes_per_sec", rate,
		"grace_period", grace.String())

	// Runtime controls so the leak can be stopped or resumed mid-session
	// without a redeploy, e.g. to stabilise the pod if a recording take goes
	// wrong.
	http.HandleFunc("/leak/stop", l.handleStop)
	http.HandleFunc("/leak/start", l.handleStart)
	http.HandleFunc("/leak/status", l.handleStatus)

	go l.run()
}

// run performs the leak: after the startup grace period it allocates and
// touches memory at a steady cadence, retaining every allocation so the pages
// stay resident until the kubelet OOMKills the container.
func (l *memoryLeaker) run() {
	if l.gracePeriod > 0 {
		time.Sleep(l.gracePeriod)
	}
	l.active.Store(true)

	// Bytes to allocate per tick, derived from the configured rate so that the
	// long-run growth matches rateBytesPerSec regardless of leakInterval.
	perTick := int(float64(l.rateBytesPerSec) * leakInterval.Seconds())
	if perTick < leakPageSize {
		perTick = leakPageSize
	}

	ticker := time.NewTicker(leakInterval)
	defer ticker.Stop()
	for range ticker.C {
		if !l.active.Load() {
			// Paused via /leak/stop: memory plateaus and the pod stays alive.
			continue
		}
		l.allocate(perTick)
	}
}

// allocate grabs n bytes, touches every page so they become resident, and
// retains the reference so the memory is never reclaimed.
func (l *memoryLeaker) allocate(n int) {
	b := make([]byte, n)
	// Write a non-zero byte to every page (and the final byte) so each page is
	// faulted in and counts towards the resident/working set. A non-zero value
	// avoids any zero-page deduplication.
	for i := 0; i < n; i += leakPageSize {
		b[i] = 0xA5
	}
	b[n-1] = 0xA5

	l.mu.Lock()
	l.blocks = append(l.blocks, b)
	l.leaked += int64(n)
	leaked := l.leaked
	l.mu.Unlock()

	l.gauge.Set(float64(leaked))
}

func (l *memoryLeaker) handleStop(w http.ResponseWriter, r *http.Request) {
	if l.active.Swap(false) {
		slog.Warn("Memory leak paused via control endpoint", "remote_addr", r.RemoteAddr)
	}
	l.writeStatus(w)
}

func (l *memoryLeaker) handleStart(w http.ResponseWriter, r *http.Request) {
	if !l.active.Swap(true) {
		slog.Warn("Memory leak resumed via control endpoint", "remote_addr", r.RemoteAddr)
	}
	l.writeStatus(w)
}

func (l *memoryLeaker) handleStatus(w http.ResponseWriter, _ *http.Request) {
	l.writeStatus(w)
}

func (l *memoryLeaker) writeStatus(w http.ResponseWriter) {
	l.mu.Lock()
	leaked := l.leaked
	l.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"active":             l.active.Load(),
		"leaked_bytes":       leaked,
		"rate_bytes_per_sec": l.rateBytesPerSec,
		"grace_period":       l.gracePeriod.String(),
	}); err != nil {
		slog.Error("Error encoding leak status response", "error", err)
	}
}
