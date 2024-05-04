package main

import (
	"log"
	"net/http"

	"github.com/butbkadrug/reaper-lyrics/pkg/projects"
)

func(app *app) getProjectHandler(w http.ResponseWriter, r *http.Request) {

    var(
        proj projects.Project
        err error
    )
    r.ParseForm()

    if id := r.FormValue("id"); id != "" {
        proj, err = app.library.GetProject(id)

        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
    }

    log.Println("rendering project: ", proj)

    err = app.renderer.RenderComponent(w, proj, "song-master")


    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
    }
}

