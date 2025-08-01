package page

import "github.com/gorilla/mux"

// AddRoutes adds the page routes to the given router
func AddRoutes(r *mux.Router) {
	r.HandleFunc("/", renderHome).Methods("GET")
	// r.HandleFunc("/blog", renderBlog).Methods("GET")
	r.HandleFunc("/{page}", renderPage).Methods("GET")
	// r.HandleFunc("/blog/articles/{article}", renderArticle).Methods("GET")
	r.HandleFunc("/projects/{project}", renderProject).Methods("GET")
}