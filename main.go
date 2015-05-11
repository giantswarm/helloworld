package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html)
	})

	go func() {
		log.Println("Starting up at :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Handle SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	log.Printf("Received signal '%v'. Exiting.", <-ch)
}

var html = `
<!DOCTYPE html>
<html lang="en">
<style>
* {
	font-family: "Helvetica Neue", "Helvetica", Helvetica, Arial, sans-serif;
}
</style>
<body>
<img src="https://giantswarm.io/static/img/logo_simplified.svg" width="200" alt="Giant Swarm">
<h1>Hello world!</h1>
<p>Congratulations you've just deployed and run your first Giant Swarm container! While you are at it, why don't you let the <a href="http://www.twitter.com/share?text=I've just deployed my first container on @giantswarm:"/>world know</a>?</p>
<p>Now go ahead and build your own Giant Swarm app in your own <a href="http://docs.giantswarm.io/guides/your-first-application/">favourite language</a>!</p>
</body>
</html>
`
