package view

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/logger"
	"html/template"
	"net/http"
)

func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("./web/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		logger.Error(err, "Error occurred while trying to generate html view.\n")
		return
	}
}
