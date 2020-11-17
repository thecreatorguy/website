package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newURL := *r.URL
		newURL.Scheme = "https"
		http.Redirect(w, r, newURL.String(), http.StatusPermanentRedirect)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}