package main

import (
	"html/template"

	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const nbsp = '\u00A0'

var headerLinks = map[string]string{
    "/":         "About",
    "/blog":     "Blog",
    "/projects": "Projects",
    "/contact":  "Contact" + string(nbsp) + "Me",
}

func multiGlobTemplate(patterns []string) *template.Template {
	tmpl := template.New("Template")
	for _, pattern := range patterns {
		tmpl = template.Must(tmpl.ParseGlob(pattern))
	}
	return tmpl
}

func main() {
	router := mux.NewRouter()
	router.Use(LogRequestMiddleware)
	router.HandleFunc("/", Home)
	assetsPath := "/assets/"
	assetsHandler := http.StripPrefix(assetsPath, http.FileServer(http.Dir("./assets/")))
	router.PathPrefix(assetsPath).Handler(assetsHandler)

	server := &http.Server{
		Addr:           ":80",
		Handler:        router,
		TLSConfig:      nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Fatal(server.ListenAndServe())
}

func LogRequestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare("/ping", r.URL.Path) != 0 {
			logrus.Info(r.Method + " " + r.URL.Path)
		}
		h.ServeHTTP(w, r)
	})
}

type BaseTemplate struct {
	URI string
	Title string
	CSSFile string
	HeaderLinks map[string]string	
	Scripts []string
}

func Home(writer http.ResponseWriter, request *http.Request) {
	var err error
	tmpl := multiGlobTemplate([]string{"./views/*.go.html", "./views/*/*.go.html"})
	err = tmpl.ExecuteTemplate(writer, "home", BaseTemplate{
		URI: "/", 
		Title: "Home",
		CSSFile: "home",
		HeaderLinks: headerLinks,
	})
	if err != nil {
		panic(err)
	}
}
