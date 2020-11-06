package page

import "github.com/gorilla/mux"

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/", home)
	
	r.HandleFunc("/{page}", page)
	r.HandleFunc("/blog/articles/{article}", article)
	r.HandleFunc("/projects/{page}", page)
}