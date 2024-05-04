package main

import (
	"net/http"

	"github.com/butbkadrug/reaper-lyrics/pkg/projects"
)

func(app *app) screenmasterHandler(w http.ResponseWriter, r *http.Request) {

    var proj projects.Project
    var projs []projects.Project

    var err error

    id := r.PathValue("id")

    param := r.URL.Query()

    if id != "" {
        proj, err = app.library.GetProject(id)

        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
    }


    if q, ok := param["search"]; ok && id != "" {
        projs = app.library.SearchProject(q[0])
    }


    type d struct {
        Project projects.Project
        Library []projects.Project
    }


    data := d {
        Project: proj,
        Library: projs,
    }

    page := Page{
        Page: Meta{
            Title: "Screen control page",
        },
        Data: data,
    }

    switch r.Method {
    case "POST":
        app.renderer.Render(w, data, "master.html")
    default:
        err = app.renderer.RenderPage(w,  page, "master.html")
    }



    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
