package view

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/logger"
	"html/template"
	"net/http"
	"time"
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

func GenerateTimeHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("./web/templates/%s.html", file))
	}
	funcMap := template.FuncMap{"fdate": formatDate}
	t := template.New("tmpl.html").Funcs(funcMap)

	templates := template.Must(t.ParseFiles(files...))
	err := templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		logger.Error(err, "Error occurred while trying to generate html view.\n")
		return
	}
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
