package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	requestNr int64  = 0
	host      string = "unknown"
)

func main() {

	currentHost, err := os.Hostname()

	if err != nil {
		log.Println("Could not determine host name:", err)
	} else {
		host = currentHost
	}

	log.Println("Starting Hello Server on " + host)

	http.HandleFunc("/", Hello)
	http.ListenAndServe(":80", nil)
}

// Hello gives out a simple hello message
func Hello(w http.ResponseWriter, r *http.Request) {

	requestNr++
	message := fmt.Sprintf("Go Hello %d from %s on %s ./%s\n", requestNr, host, r.Method, r.URL.Path[1:])
	log.Print(message)
	fmt.Fprint(w, message)
}
