// Package app creates the web server with all routes and middleware.
package app

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"website/internal/app/message"
	"website/internal/app/page"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	cards "github.com/thecreatorguy/cards/pkg/web"
	"github.com/thecreatorguy/shakesearch/pkg/shakesearch"
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

	// Shakesearch project import
	searcher := shakesearch.Searcher{}
	logrus.Info("Loading in the works of Shakespeare!..")
	err := searcher.Load("./data/completeworks_shakespeare.txt")
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("Loaded!")
	shakesearch.AddRoutes(r, searcher, 
		"./shakesearch/views/index.go.html", 
		"./shakesearch/static", 
		"/projects/shakesearch", 
		"/app",
		"/assets",
	)
	
	cards.AddRoutes(r, "/projects/cards")
	
	server := &http.Server{
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if os.Getenv("HTTPS") != "1" {
		server.Addr = ":8675"
		logrus.Fatal(server.ListenAndServe())
		logrus.Info("Server ready to handle requests!")
	} else {
		cert, err := tls.LoadX509KeyPair(os.Getenv("SSL_CERT_PATH"), os.Getenv("SSL_KEYFILE_PATH"))
		if err != nil {
			logrus.WithError(err).Fatal("Couldn't create ssl certificate")
		}
		server.Addr = ":8676"
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		logrus.Fatal(server.ListenAndServeTLS(os.Getenv("SSL_CERT_PATH"), os.Getenv("SSL_KEYFILE_PATH")))
		logrus.Info("Server ready to handle requests!")
	}	
}
