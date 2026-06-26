package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var gitCommit = "n/a"

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
	http.Handle("/echo", loggingHandler(http.HandlerFunc(echoHandler)))

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

// Healthz endpoint
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, "OK\n")
	if err != nil {
		slog.Error("Error in io.WriteString", "error", err)
	}
}

// echoResponse describes the request echoed back by echoHandler.
type echoResponse struct {
	ReceivedAt  time.Time           `json:"received_at"`
	RespondedAt time.Time           `json:"responded_at"`
	Delay       string              `json:"delay"`
	Method      string              `json:"method"`
	Path        string              `json:"path"`
	Query       map[string][]string `json:"query"`
	Host        string              `json:"host"`
	RemoteAddr  string              `json:"remote_addr"`
	Proto       string              `json:"proto"`
	Headers     map[string][]string `json:"headers"`
	Body        string              `json:"body"`
}

// Echo endpoint: echoes back the request details as JSON, optionally delaying
// the response by the duration given in the "delay" query param (e.g. ?delay=2s).
func echoHandler(w http.ResponseWriter, r *http.Request) {
	receivedAt := time.Now().UTC()

	var delay time.Duration
	if raw := r.URL.Query().Get("delay"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil || d < 0 {
			http.Error(w, "invalid delay: expected a non-negative Go duration like 500ms, 2s or 1m", http.StatusBadRequest)
			return
		}
		delay = d
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading request body", "error", err)
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-r.Context().Done():
			slog.Info("Echo request cancelled during delay", "error", r.Context().Err())
			return
		}
	}

	resp := echoResponse{
		ReceivedAt:  receivedAt,
		RespondedAt: time.Now().UTC(),
		Delay:       delay.String(),
		Method:      r.Method,
		Path:        r.URL.Path,
		Query:       r.URL.Query(),
		Host:        r.Host,
		RemoteAddr:  r.RemoteAddr,
		Proto:       r.Proto,
		Headers:     r.Header,
		Body:        string(body),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Error encoding echo response", "error", err)
	}
}
