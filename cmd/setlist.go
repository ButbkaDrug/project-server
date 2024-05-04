package main

import (
	"net/http"
)

func (app *app) setlistHandler(w http.ResponseWriter, r *http.Request) {



    var err error
	file := "setlist.html"
    data := app.setlist.GetSongs()

    page := Page {
        Page: Meta{
            Title: "Setlist",
        },
        Data: data,
    }

    switch r.Method {

    case "POST":
        err = app.renderer.RenderComponent(w, data, "set")
    default:
        err = app.renderer.RenderPage(w, page, file)

    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

}
