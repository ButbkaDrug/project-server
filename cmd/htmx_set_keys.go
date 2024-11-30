package main

import (
	"fmt"
	"net/http"
)

func (app *app) htmxSetKeysHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	keys := app.ParseKeysArgument(&r.Form)

	url := "/guitar/?"

	for i, key := range keys {
		if key == "" {
			continue
		}

		url += fmt.Sprintf("key%d=%s&", i+1, key)
	}

	w.Header().Set("HX-Replace-Url", url)

	track_items := app.reaper.Project.Tracks[0].Items
	items := make([]MapItem, len(track_items))

	for i, item := range track_items {
		items[i].ItemID = i + 1
		items[i].Name = item.Name[:20]
	}

	builder := StagePageBuilder{}
	builder.SetRenderer(app.renderer).
		SetItems(items).
		SetData(app.reaper.Sections).
		SetKeys(keys)

	master := PageMaster{}

	page, err := master.CreateChordedPage(&builder)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	err = app.renderer.RenderComponent(w, page.Data, "scroll")

	if err != nil {
		app.logger.Error("failed to execute template: ", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
