package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "https"
		r.URL.Host = "itstimjohnson.com"
		http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}