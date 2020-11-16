// Package page implements routes serving the web pages for the itstimjohnson website.
package page

import (
	"bytes"
	"html/template"
	"net/http"
	"website/internal/app/response"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type pageInput struct {
	URI string
	Title string
	CSSFile string
	JSScripts []template.HTMLAttr
	JSONData map[string]template.JS
	PageTemplateName string
	PageTemplateData interface{}
}

func (p pageInput) RenderPage() template.HTML {
	var buf bytes.Buffer
	err := pageTemplates.ExecuteTemplate(&buf, p.PageTemplateName, p.PageTemplateData)
	if err != nil {
		panic(err)
	}

	return template.HTML(buf.Bytes())
}

var topLevelTemplates, pageTemplates *template.Template

// init initializes the templates so they only have to be read once
func init() {
	topLevelTemplates = template.Must(template.ParseGlob("./views/*.go.html"))
	pageTemplates = template.Must(template.ParseGlob("./views/pages/*.go.html"))
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
		JSScripts: []template.HTMLAttr{"src=\"https://www.google.com/recaptcha/api.js\" async defer"},
	},
	"confirmation": {
		Title: "Confirmation",
		CSSFile: "confirmation",
	},
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	if _, found := pageNameToInput[page]; !found {
		response.Write404(w)
		return
	}

	input := pageNameToInput[page]
	input.URI = r.URL.Path
	input.PageTemplateName = page
	render(w, "page", input)
}

func render(w http.ResponseWriter, template string, data pageInput) {
	var buf bytes.Buffer
	err := topLevelTemplates.ExecuteTemplate(&buf, template, data)
	if err != nil {
		logrus.WithError(err).Error("Failed executing template")
		response.Write500(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
