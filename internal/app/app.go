// Package app creates the web server with all routes and middleware.
package app

import (
	"net/http"
	"strings"
	"time"
	"website/internal/app/message"
	"website/internal/app/page"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	cards "github.com/thecreatorguy/cards/pkg/web"
)

// StartWebServer starts the webserver for the website
func StartWebServer() {
	r := mux.NewRouter()

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Compare("/ping", r.URL.Path) != 0 {
				logrus.Info(r.Method + " " + r.URL.Path)
			}
			h.ServeHTTP(w, r)
		})
	})

	page.AddRoutes(r)
	message.AddRoutes(r)
	addAPIRoutes(r)

	assetsPath := "/assets/"
	assetsHandler := http.StripPrefix(assetsPath, http.FileServer(http.Dir("./assets/")))
	r.PathPrefix(assetsPath).Handler(assetsHandler)
	
	cards.AddRoutes(r, "/projects/cards", "/assets/images")
	
	server := &http.Server{
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Addr: ":8675",
	}
	
	logrus.Info("Starting server...")
	logrus.Fatal(server.ListenAndServe())
}
