package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (app *app) autoscrollHandler(w http.ResponseWriter, r *http.Request) {
	file := "autoscroll.html"

	track_items := app.reaper.Project.Tracks[0].Items
	items := make([]MapItem, len(track_items))

	for i, item := range track_items {
		items[i].ItemID = i + 1
		items[i].Name, _, _ = strings.Cut(item.Name, "-")

		for _, line := range app.reaper.Lines {
			if line.Position == item.Position {
				items[i].LineID = line.ID
			}
		}
	}

	builder := StagePageBuilder{}
	builder.SetData(app.reaper.Sections).SetRenderer(app.renderer)
	builder.SetItems(items)

	master := PageMaster{}
	page, err := master.CreateSingersPage(&builder)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if err = app.renderer.RenderPage(w, page, file); err != nil {

		app.logger.Error("failed to execute template: ", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
