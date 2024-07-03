package main

import (
	"net/http"
)

func (app *app) autoscrollHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	file := "autoscroll.html"
	params := r.URL.Query()

	page := Page{
		Page: Meta{
			Title: "Lyrics will scroll automatically",
		},
		Map:  nil,
		Data: app.reaper.Sections,
	}
	if params.Has("map") {

		page.Map = app.reaper.GetMapForCurrentItem()
	}

	switch r.Method {
	case "POST":
		err = app.renderer.RenderComponent(w, page.Data, "scroll")
	default:
		err = app.renderer.RenderPage(w, page, file)

	}

	if err != nil {
		app.logger.Error("failed to execute template: ", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
