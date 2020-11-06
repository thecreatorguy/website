package page

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"

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

func init() {
	TopLevelTemplates = template.Must(template.ParseGlob("./views/*.go.html"))
	PageTemplates = template.Must(template.ParseGlob("./views/pages/*.go.html"))

	temp, err := ioutil.ReadFile("./data/slider-levels.json")
	if err != nil {
		panic(err)
	}
	sliderLevels = template.JS(temp)
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
