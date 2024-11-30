package main

import (
	"fmt"
	"net/http"
)

func (app *app) drumsHandler(w http.ResponseWriter, r *http.Request) {

	templateFile := "autoscroll.html"

	builder := StagePageBuilder{}

	page, err := builder.SetRenderer(app.renderer).
		SetTitle("No nonsense drum page").
		SetDescription("You not gonna see silly lyrics or chords here. Only big comment, only hard core!").
		SetData(app.reaper.Sections).
		SetComments("drums").
		WithBigComments().
		Build()

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	err = app.renderer.RenderPage(w, page, templateFile)

	if err != nil {
		app.logger.Error("failed to execute template: ", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
