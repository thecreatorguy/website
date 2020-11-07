package message

import "github.com/gorilla/mux"

// AddRoutes adds the message routes to the given router
func AddRoutes(r *mux.Router) {
	r.HandleFunc("/message", sendMessageEndpoint).Methods("POST")
}