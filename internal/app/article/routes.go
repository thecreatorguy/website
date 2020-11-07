package article

import "github.com/gorilla/mux"

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/articles/{article}", getArticleEndpoint).Methods("GET")
	r.HandleFunc("/articles/", createArticleEndpoint).Methods("POST")
	r.HandleFunc("/articles/{article}", updateArticleEndpoint).Methods("PUT")
	r.HandleFunc("/articles/{article}", deleteArticleEndpoint).Methods("DELETE")
}