package main

import (
	"net/http"
)

func(app *app) addtosetlistHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    r.ParseForm()

    v := r.FormValue("path")

    if app.setlist.Contains(v) {
        w.WriteHeader(http.StatusUnprocessableEntity)
        app.renderer.Render(w, "song already in the setlist", "alert")
        return
    }

    song, err := app.setlist.Append(v)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    app.renderer.RenderComponent(w, song, "oob-set-element")
}
