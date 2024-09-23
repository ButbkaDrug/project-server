package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *app) guitarHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello from guitar route")
	templateFile := "autoscroll.html"
	params := r.URL.Query()

	capo := app.ParseCapoArgument(params.Get("capo"))
	keys := app.ParseKeysArgument(params.Get("keys"))

	htmlString := RenderContent(
		&RenderContentConfig{
			Data:         app.reaper.Sections,
			Render:       app.renderer,
			CommentsType: "Guitar",
			WithChords:   true,
			WithLyrics:   true,
			Keys:         keys,
			Capo:         capo,
		},
	)

	page := Page{
		Page: Meta{
			Title: "Lyrics will scroll automatically",
			Stats: []stat{
				{
					Title: "IP address of the server",
					Value: fmt.Sprintf("%s%s", app.IP, app.Port),
					Desc:  "Make sure you are on the same local network and use this address",
				},
			},
		},
		Map: nil,
		Data: struct {
			Content interface{}
		}{
			template.HTML(htmlString),
		},
	}

	var err error

	switch r.Method {
	case "POST":
		err = app.renderer.RenderComponent(w, page.Data, "scroll")
	default:
		err = app.renderer.RenderPage(w, page, templateFile)

	}

	if err != nil {
		app.logger.Error("failed to execute template: ", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
