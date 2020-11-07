package app

import (
	"net/http"
	"website/internal/app/article"

	"github.com/gorilla/mux"
)

func addAPIRoutes(r *mux.Router) {
	s := r.PathPrefix("/api").Subrouter()
	s.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			h.ServeHTTP(w, r)
		})
	})
	
	article.AddRoutes(s)
}