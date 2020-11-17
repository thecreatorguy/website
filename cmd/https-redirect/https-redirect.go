package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// r.URL.Scheme = "https"
		fmt.Printf(r.URL.String())
		http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}