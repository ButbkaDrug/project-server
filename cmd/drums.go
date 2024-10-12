package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *app) drumsHandler(w http.ResponseWriter, r *http.Request) {

	templateFile := "autoscroll.html"
	params := r.URL.Query()

	var withChords bool
	var withLyrics bool
	var commentsType string = "drums"

	if c, err := strconv.ParseBool(params.Get("chords")); err == nil {
		withChords = c
	}

	if l, err := strconv.ParseBool(params.Get("lyrics")); err == nil {
		withLyrics = l
	}

	if ct := params.Get("comments"); ct != "" {
		commentsType = ct
	}

	keys := app.ParseKeysArgument(&params)

	htmlString := RenderContent(
		&RenderContentConfig{
			logger:       app.logger,
			Data:         app.reaper.Sections,
			Render:       app.renderer,
			CommentsType: commentsType,
			CommentsBig:  true,
			WithChords:   withChords,
			WithLyrics:   withLyrics,
			Keys:         keys,
		},
	)

	page := Page{
		Page: Meta{
			Title: "Drummer page has no lyrics or chords, but have big comments instead",
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
