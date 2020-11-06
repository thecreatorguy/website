package page

import (
	"html/template"
	"io/ioutil"
)

func readJSONFile(filename string) template.JS {
	temp, err := ioutil.ReadFile("./data/slider-levels.json")
	if err != nil {
		panic(err)
	}
	return template.JS(temp)
}