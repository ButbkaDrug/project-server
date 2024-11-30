package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (app *app) guitarHandler(w http.ResponseWriter, r *http.Request) {
	templateFile := "autoscroll.html"
	params := r.URL.Query()

	keys := app.ParseKeysArgument(&params)

	track_items := app.reaper.Project.Tracks[0].Items
	items := make([]MapItem, len(track_items))

	for i, item := range track_items {
		items[i].ItemID = i + 1
		items[i].Name, _, _ = strings.Cut(item.Name, "-")

		for _, line := range app.reaper.Lines {
			fmt.Println(line.Position, item.Position)

			if line.Position == item.Position {
				items[i].LineID = line.ID
			}
		}
	}
	builder := StagePageBuilder{}
	builder.SetRenderer(app.renderer).
		SetData(app.reaper.Sections).
		SetItems(items).
		SetKeys(keys)
	master := PageMaster{}

	page, err := master.CreateChordedPage(&builder)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if err = app.renderer.RenderPage(w, page, templateFile); err != nil {
		app.logger.Error("failed to execute template: ", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
