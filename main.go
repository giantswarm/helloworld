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

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

var html = `
<!DOCTYPE html>
<html lang="en">
<body>

<h1>Hello world!</h1>

<p>Congratulations you've just deployed and run your first Giant Swarm container!</p> 
<p>Now go ahead and build your own Giant Swarm app with your own <a href="http://docs.giantswarm.io/guides/your-first-application/">favourite language!</a></p>.

</body>
</html>
`
