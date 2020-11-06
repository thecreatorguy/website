package page

import (
	"bytes"
	"html/template"
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

type pageInput struct {
	URI string
	Title string
	CSSFile string
	JSScripts []string
	JSONData map[string]template.JS
	PageTemplateName string
	PageTemplateData string
}

func (p pageInput) renderPage() template.HTML {

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

func renderPage(writer http.ResponseWriter, request *http.Request) {
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

	render(writer, "page", PageTemplate{
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
