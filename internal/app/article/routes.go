package article

import "github.com/gorilla/mux"

// AddRoutes adds the api routes for manipulating articles in the database to the given router
func AddRoutes(r *mux.Router) {
	r.HandleFunc("/articles", getArticlesEndpoint).Methods("GET")
	r.HandleFunc("/articles/{article}", getArticleEndpoint).Methods("GET")
	r.HandleFunc("/articles", createArticleEndpoint).Methods("POST")
	r.HandleFunc("/articles/{article}", updateArticleEndpoint).Methods("PUT")
	r.HandleFunc("/articles/{article}", deleteArticleEndpoint).Methods("DELETE")
}