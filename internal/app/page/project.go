package page

import (
	"html/template"
	"io/ioutil"
	"net/http"

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
		JSScripts: []string{"slider"},
		JSONData: map[string]template.JS{"level-data": readJSONFile("./data/slider-levels.json")},
	},
	"jumpybird": {
		Title: "Jumpy Bird AI",
		CSSFile: "jumpybird",
		JSScripts: []string{"jumpybird"},
	},
}


func renderProject(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	if _, found := projectNameToInput[project]; !found {
		write404(w)
		return
	}

	input := projectNameToInput[project]
	input.URI = r.URL.Path
	input.PageTemplateName = project
	render(w, "page", input)
}