package app

import (
	"net/http"
	"os"
	"website/internal/app/article"
	"website/internal/app/response"

	"github.com/gorilla/mux"
)

func addAPIRoutes(router *mux.Router) {
	s := router.PathPrefix("/api").Subrouter()
	s.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-API-Key") != os.Getenv("API_KEY") {
				response.Write401(w)
				return
			}

			h.ServeHTTP(w, r)
		})
	})
	
	article.AddRoutes(s)
}