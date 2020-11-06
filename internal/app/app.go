package app

import (
	"net/http"
	"strings"
	"time"
	"website/internal/app/article"
	"website/internal/app/page"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func StartWebServer() {
	router := mux.NewRouter()

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Compare("/ping", r.URL.Path) != 0 {
				logrus.Info(r.Method + " " + r.URL.Path)
			}
			h.ServeHTTP(w, r)
		})
	})

	page.AddRoutes(router)
	article.AddRoutes(router)

	assetsPath := "/assets/"
	assetsHandler := http.StripPrefix(assetsPath, http.FileServer(http.Dir("./assets/")))
	router.PathPrefix(assetsPath).Handler(assetsHandler)

	
	server := &http.Server{
		Addr:           ":8675",
		Handler:        router,
		TLSConfig:      nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Fatal(server.ListenAndServe())
}
