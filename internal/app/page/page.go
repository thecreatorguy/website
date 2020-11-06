package page

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type pageInput struct {
	URI string
	Title string
	CSSFile string
	JSScripts []string
	JSONData map[string]template.JS
	PageTemplateName string
	PageTemplateData string
}

func (p pageInput) RenderPage() template.HTML {
	var buf bytes.Buffer
	err := PageTemplates.ExecuteTemplate(&buf, p.PageTemplateName, p.PageTemplateData)
	if err != nil {
		panic(err)
	}

	return template.HTML(buf.Bytes())
}

var TopLevelTemplates, PageTemplates *template.Template

func init() {
	TopLevelTemplates = template.Must(template.ParseGlob("./views/*.go.html"))
	PageTemplates = template.Must(template.ParseGlob("./views/pages/*.go.html"))
}

func renderHome(w http.ResponseWriter, request *http.Request) {
	render(w, "home", pageInput{
		URI: "/", 
		Title: "Home",
		CSSFile: "home",
	})
}

var pageNameToInput = map[string]pageInput{
	"projects": {
		Title: "Projects",
		CSSFile: "projects",
	},
	"contact": {
		Title: "Contact Me",
		CSSFile: "contact",
	},
	"confirmation": {
		Title: "Confirmation",
		CSSFile: "confirmation",
	},
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	if _, found := pageNameToInput[page]; !found {
		write404(w)
		return
	}

	input := pageNameToInput[page]
	input.URI = r.URL.Path
	input.PageTemplateName = page
	render(w, "page", input)
}

func render(w http.ResponseWriter, template string, data interface{}) {
	var buf bytes.Buffer
	err := TopLevelTemplates.ExecuteTemplate(&buf, template, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Failed executing template")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func write404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}