package message

import "github.com/gorilla/mux"

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/message", sendMessageEndpoint).Methods("POST")
}