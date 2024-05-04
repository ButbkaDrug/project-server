package main

import (
	"log"
	"net/http"
)

func (app *app) viewsongHandler(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")
	filename := "single-project.html"

	proj, err := app.library.GetProject(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	page := Page{
		Page: Meta{
			Title: proj.Name,
		},
		Data: proj,
	}

	log.Println(proj)

	switch r.Method {
	case "POST":
		err = app.renderer.Render(w, proj, filename)
	default:
		err = app.renderer.RenderPage(w, page, filename)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
