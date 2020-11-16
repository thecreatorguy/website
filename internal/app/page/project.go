package page

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"website/internal/app/response"

	"github.com/gorilla/mux"
)

func readJSONFile(filename string) template.JS {
	temp, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return template.JS(temp)
}

var projectNameToInput = map[string]pageInput{
	"slider": {
		Title: "Slider Game",
		CSSFile: "slider",
		JSScripts: []template.HTMLAttr{"src=\"assets/js/slider.js\""},
		JSONData: map[string]template.JS{"level-data": readJSONFile("./data/slider-levels.json")},
	},
	"jumpybird": {
		Title: "Jumpy Bird AI",
		CSSFile: "jumpybird",
		JSScripts: []template.HTMLAttr{"src=\"assets/js/jumpybird.js\""},
	},
}


func renderProject(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	if _, found := projectNameToInput[project]; !found {
		response.Write404(w)
		return
	}

	input := projectNameToInput[project]
	input.URI = r.URL.Path
	input.PageTemplateName = project
	render(w, "page", input)
}