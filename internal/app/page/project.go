package page

import (
	"fmt"
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
		JSScripts: []template.HTMLAttr{"src=\"/assets/js/slider.js\""},
		JSONData: map[string]template.JS{"level-data": readJSONFile("./data/slider-levels.json")},
	},
	"jumpybird": {
		Title: "Jumpy Bird AI",
		CSSFile: "jumpybird",
		JSScripts: []template.HTMLAttr{"src=\"/assets/js/jumpybird.js\""},
	},
}

var projectNameToRedirect = map[string]string{
	"cards": "waitingroom",
}


func renderProject(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	if input, ok := projectNameToInput[project]; ok {
		input.URI = r.URL.Path
		input.PageTemplateName = project
		render(w, "page", input)
		return
	}
	if redirect, ok := projectNameToRedirect[project]; ok {
		http.Redirect(w, r, fmt.Sprintf("/projects/%s/%s", project, redirect), http.StatusSeeOther)
		return
	}

	response.Write404(w)
}