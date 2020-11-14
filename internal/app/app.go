// Package app creates the web server with all routes and middleware.
package app

import (
	"crypto/tls"
	"net/http"
	"os"
	"strings"
	"time"
	"website/internal/app/message"
	"website/internal/app/page"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	server := &http.Server{
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if strings.ToLower(os.Getenv("HTTPS")) != "true" {
		server.Addr = ":8675"
		logrus.Info("Server ready to handle requests!")
		logrus.Fatal(server.ListenAndServe())
	} else {
		kpr, err := NewKeypairReloader("", "")
		if err != nil {
			logrus.WithError(err).Fatal("Couldn't create key pair loader")
		}
		server.Addr = ":8676"
		server.TLSConfig = &tls.Config{
			GetCertificate: kpr.GetCertificateFunc(),
		}
		logrus.Info("Server ready to handle requests!")
		logrus.Fatal(server.ListenAndServeTLS("", ""))
	}	
}
