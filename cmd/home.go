package main

import (
	"log/slog"
	"net/http"
)

type TableOfContentEntry struct {
	Title string
	Link  string
}

type Meta struct {
	Title          string
	TableOfContent []TableOfContentEntry
}

type Page struct {
	Page Meta
	Data interface{}
}

func (app *app) noscrollHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	data := []Section{}

	for _, item := range app.reaper.Items {
		data = append(data, Section{
			ID:    len(data),
			Name:  item.Name,
			Type:  TypeHero,
			Image: item.Image,
		})

		for _, section := range item.Lyrics.Sections {

			data = append(data, Section{

				ID:   len(data),
				Name: section.Name,
				Type: TypeSection,
				Text: section.Text,
			})

		}

	}

	page := Page{
		Page: Meta{
			Title: "Welcome to the church!",
		},
		Data: data,
	}

	switch r.Method {
	case "POST":

		err = app.renderer.Render(w, data, "index.html", "components.html")

	default:
		err = app.renderer.RenderPage(w, page, "index.html")
	}

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
