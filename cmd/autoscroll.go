package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *app) autoscrollHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cmd/autoscroll.go invoked")
	var err error

	file := "autoscroll.html"
	params := r.URL.Query()

	displayChords := false

	if c, err := strconv.ParseBool(params.Get("chords")); err == nil {
		displayChords = c
	}

	displayLyrics := true

	if l, err := strconv.ParseBool(params.Get("lyrics")); err == nil {
		displayLyrics = l
	}

	cms := params.Get("comments")

	keys := app.ParseKeysArgument(params.Get("keys"))
	capo := app.ParseCapoArgument(params.Get("capo"))

	fmt.Println("about to render content")
	htmlString := RenderContent(
		&RenderContentConfig{
			Data:         app.reaper.Sections,
			Render:       app.renderer,
			CommentsType: cms,
			WithChords:   displayChords,
			WithLyrics:   displayLyrics,
			Keys:         keys,
			Capo:         capo,
		},
	)

	app.logger.Debug("rendering autoscroll", "chords", displayChords, "lyrics", displayLyrics)

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

	switch r.Method {
	case "POST":
		err = app.renderer.RenderComponent(w, page.Data, "scroll")
	default:
		err = app.renderer.RenderPage(w, page, file)

	}

	if err != nil {
		app.logger.Error("failed to execute template: ", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
