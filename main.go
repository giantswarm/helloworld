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
	"syscall"

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
