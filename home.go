package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

func (app *app) homeHandler(w http.ResponseWriter, r *http.Request) {

	f := path.Join("pkg", "templates", "main.gotmpl")

	t, err := template.ParseFS(templateMain, f)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, app.reaper)

	if err != nil {
		log.Println("failed to execute template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
