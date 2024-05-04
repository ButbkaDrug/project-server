package main

import (
	"net/http"
)

func(app *app) projectsearchHandler(w http.ResponseWriter, r *http.Request) {

    query := r.FormValue("search")

    app.renderer.RenderComponent(w, app.library.SearchProject(query), "project-search-results")
}
