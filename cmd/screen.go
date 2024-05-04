package main

import (
	"net/http"
)

func(app *app) screenHandler(w http.ResponseWriter, r *http.Request) {

    app.renderer.Render(w, app.screen.Data, "screen.html")

}
