package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func(app *app) updateSetlist(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    r.ParseForm()

    for key, val := range r.PostForm {

        switch key {
        case "delete":
            for _, v := range val {
                id, err := strconv.Atoi(v)

                if err != nil {
                    http.Error(w, "wrong value type", http.StatusInternalServerError)
                    return
                }

                err = app.setlist.Delete(id)

                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }

            }

        case "up":

            for _, v := range val {
                id, err := strconv.Atoi(v)

                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }

                err = app.setlist.MoveUp(id)

                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return

                }
            }

        case "down":

            for _, v := range val {

                id, err := strconv.Atoi(v)

                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }

                err = app.setlist.MoveDown(id)

                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return

                }
            }

        default:
            fmt.Fprintln(w, val)
        }
    }

    app.renderer.RenderComponent(w, app.setlist.GetSongs(), "set")
}
