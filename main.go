package main

import (
	"fmt"
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
		if strings.ToLower(v) == "version" {
			fmt.Printf("helloworld version: %s\n", gitCommit)
			os.Exit(0)
		}
	}

	mime.AddExtensionType(".ico", "image/x-icon")
	mime.AddExtensionType(".svg", "image/svg+xml")

	fileHandler := http.FileServer(http.Dir("/content"))
	http.Handle("/", loggingHandler(fileHandler))

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
