package main

import (
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var gitCommit = "n/a"

// heartbeatOK tracks whether the last cronitor heartbeat ping succeeded.
// The app is considered not ready if this is false.
var heartbeatOK atomic.Bool

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"code"},
	)
)

func main() {
	for _, v := range os.Args {
		param := strings.ToLower(v)
		switch param {
		case "version":
			fmt.Printf("helloworld version: %s\n", gitCommit)
			os.Exit(0)
		case "--help":
			fmt.Printf("usage: %s\n", os.Args[0])
			os.Exit(0)
		}

	}

	// Read SECRET_KEY environment variable and fail if not set
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		slog.Error("SECRET_KEY environment variable is required but not set")
		os.Exit(1)
	}
	slog.Info("SECRET_KEY loaded successfully")

	err := mime.AddExtensionType(".ico", "image/x-icon")
	if err != nil {
		slog.Error("Error when adding mime type for .ico", "error", err)
	}
	err = mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		slog.Error("Error when adding mime type for .svg", "error", err)
	}

	fileHandler := http.FileServer(http.Dir("/content"))
	http.Handle("/", loggingHandler(fileHandler))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readyz", readyzHandler)

	// Start heartbeat pinger in background.
	heartbeatURL := os.Getenv("HEARTBEAT_URL")
	if heartbeatURL == "" {
		slog.Error("HEARTBEAT_URL environment variable is required but not set")
		os.Exit(1)
	}
	go pingHeartbeat(heartbeatURL)

	go func() {
		slog.Info("Starting up at :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil { //nolint:gosec
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Handle SIGTERM.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	slog.Info("Received signal, exiting", "signal", <-ch)
}

// HTTP Handler that adds logging to STDOUT
func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)
		h.ServeHTTP(w, r)
		httpRequestsTotal.With(prometheus.Labels{"code": w.Header().Get("Code")}).Inc()
	})
}

// Healthz endpoint (liveness)
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, "OK\n")
	if err != nil {
		slog.Error("Error in io.WriteString", "error", err)
	}
}

// Readyz endpoint (readiness) — fails when the heartbeat ping is unhealthy.
func readyzHandler(w http.ResponseWriter, r *http.Request) {
	if !heartbeatOK.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = io.WriteString(w, "heartbeat unhealthy\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "OK\n")
}

// pingHeartbeat pings the heartbeat URL every minute and updates heartbeatOK.
func pingHeartbeat(url string) {
	client := &http.Client{Timeout: 10 * time.Second}

	ping := func() {
		resp, err := client.Get(url) //nolint:gosec
		if err != nil {
			slog.Error("Heartbeat ping failed", "error", err)
			heartbeatOK.Store(false)
			return
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			heartbeatOK.Store(true)
		} else {
			slog.Error("Heartbeat ping returned unexpected status", "status", resp.StatusCode)
			heartbeatOK.Store(false)
		}
	}

	// Ping immediately on startup.
	ping()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		ping()
	}
}
