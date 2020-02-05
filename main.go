package main

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mime.AddExtensionType(".ico", "image/x-icon")
	mime.AddExtensionType(".svg", "image/svg+xml")

	fileHandler := http.FileServer(http.Dir("/content"))
	http.Handle("/", loggingHandler(fileHandler))
	http.HandleFunc("/healthz", healthzHandler)

	go func() {
		log.Println("Starting up at :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Handle SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	log.Printf("Received signal '%v'. Exiting.", <-ch)
}

// HTTP Handler that adds logging to STDOUT
func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

// Healthz endpoint
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK\n")
}
