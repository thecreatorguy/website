package page

import "github.com/gorilla/mux"

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/", renderHome)
	r.HandleFunc("/blog", renderBlog)
	r.HandleFunc("/{page}", renderPage)
	r.HandleFunc("/blog/articles/{article}", renderArticle)
	r.HandleFunc("/projects/{project}", renderProject)
}