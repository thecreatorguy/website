package main

import (
	"bytes"
	"html/template"
	"io/ioutil"

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

var pageTitles = map[string]string{
	"home":          "Home",
	"projects":      "Projects",
	"contact":       "Contact Me",
	"confirmation":  "Confirmation",
	"slider":        "Slider Game",
	"resume":        "Resume",
	"jumpybird":     "Jumpy Bird AI",
}

type BasePageTemplate struct {
	URI string
	Title string
	CSSFile string
	HeaderLinks map[string]string	
	Scripts []string
}

type PageTemplate struct {
	BasePageTemplate
	JSONData map[string]template.JS
	RenderedPage template.HTML
}

var TopLevelTemplates, PageTemplates *template.Template
var sliderLevels template.JS

func main() {
	TopLevelTemplates = template.Must(template.ParseGlob("./views/*.go.html"))
	PageTemplates = template.Must(template.ParseGlob("./views/pages/*.go.html"))

	temp, err := ioutil.ReadFile("./data/slider-levels.json")
	if err != nil {
		panic(err)
	}
	sliderLevels = template.JS(temp)


	router := mux.NewRouter()

	router.Use(logRequestMiddleware)

	router.HandleFunc("/", home)
	
	router.HandleFunc("/{page}", page)
	// router.HandleFunc("/blog/articles/{article}", article)
	router.HandleFunc("/projects/{page}", page)

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

func logRequestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare("/ping", r.URL.Path) != 0 {
			logrus.Info(r.Method + " " + r.URL.Path)
		}
		h.ServeHTTP(w, r)
	})
}

func home(writer http.ResponseWriter, request *http.Request) {
	basePageExecute(writer, TopLevelTemplates, "home", BasePageTemplate{
		URI: "/", 
		Title: "Home",
		CSSFile: "home",
		HeaderLinks: headerLinks,
	})
}

func page(writer http.ResponseWriter, request *http.Request) {
	page := mux.Vars(request)["page"]

	scripts := []string{}
	jsonData := map[string]template.JS{}
	if page == "jumpybird" || page == "slider" {
		scripts = append(scripts, page)
	}
	if page == "slider" {
		jsonData["level-data"] = sliderLevels
	}

	var buf bytes.Buffer
	err := PageTemplates.ExecuteTemplate(&buf, page, nil)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Failed executing template")
		return
	}

	basePageExecute(writer, TopLevelTemplates, "page", PageTemplate{
		BasePageTemplate: BasePageTemplate{
			URI: request.URL.Path, 
			Title: pageTitles[page],
			CSSFile: page,
			HeaderLinks: headerLinks,
			Scripts: scripts,
		},
		RenderedPage: template.HTML(buf.Bytes()),
		JSONData: jsonData,
	})
}

func basePageExecute(writer http.ResponseWriter, tmpl *template.Template, template string, data interface{}) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, template, data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Failed executing template")
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(buf.Bytes())
}
