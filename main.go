package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var gitCommit = "n/a"

func main() {
	for _, v := range os.Args {
		param := strings.ToLower(v)
		switch {
		case param == "version":
			fmt.Printf("helloworld version: %s\n", gitCommit)
			os.Exit(0)
		case param == "--help":
			fmt.Printf("usage: %s\n", os.Args[0])
			os.Exit(0)
		}

	}

	err := mime.AddExtensionType(".ico", "image/x-icon")
	if err != nil {
		log.Printf("Error when adding mime type for .ico: %s", err)
	}
	mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		log.Printf("Error when adding mime type for .svg: %s", err)
	}

	fileHandler := http.FileServer(http.Dir("/content"))
	http.Handle("/", loggingHandler(fileHandler))
	http.HandleFunc("/healthz", healthzHandler)

	go func() {
		log.Println("Starting up at :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Handle SIGTERM.
	ch := make(chan os.Signal, 1)
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
