package main

import (
	"fmt"
	"net/http"
)

func(app *app) updatescreenHandler(w http.ResponseWriter, r *http.Request){
    id := r.FormValue("id")
    section := r.FormValue("section")

    if id == section {
        http.Error(w, "empty request not allowed", http.StatusBadRequest)
        return
    }

    proj, err  := app.library.GetProject(id)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sec, err := proj.GetSection(section)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    go app.screen.Clients.Broadcast(sec.Text)

    fmt.Fprintf(w, "%s", "success")

}
