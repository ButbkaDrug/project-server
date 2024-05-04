package main

import (
	"log"
	"net/http"
)

func (app *app) songsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	var err error
	file := "songs.html"

	page := Page{
		Page: Meta{
			Title: "Song library",
		},

		Data: app.library.Projects,
	}

	switch r.Method {
	case "POST":

        if err = r.ParseForm(); err != nil {
            log.Println("cmd/songs.go faild to parse form: ", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }


		if query := r.FormValue("search"); len(r.Form) > 0 {
			data := app.library.SearchProject(query)
			err = app.renderer.RenderComponent(w, data, "songs-table-content")
		} else {

            err = app.renderer.RenderComponent(w, page.Data, "songs-table")
        }


	default:
		err = app.renderer.RenderPage(w, page, file)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
