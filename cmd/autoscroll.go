// TODO: I want to move all the logic from template into here and have as little logic inside the telplate as possible
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"lp/pkg/lyrics"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/butbkadrug/reaper-lyrics/pkg/reaper"
	"github.com/butbkadrug/reaper-lyrics/pkg/templates"
)

func (app *app) autoscrollHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	file := "autoscroll.html"
	params := r.URL.Query()

	displayChords := false

	app.logger.Debug("autocrollHandler", "params", params)

	if c, err := strconv.ParseBool(params.Get("chords")); err == nil {
		displayChords = c
	}

	displayLyrics := true

	if l, err := strconv.ParseBool(params.Get("lyrics")); err == nil {
		displayLyrics = l
	}

	keys := make([]string, len(app.reaper.Track.Items))
	if k := strings.Split(params.Get("keys"), ","); len(k) > 0 {
		checklist := []string{
			"A",
			"Bb",
			"B",
			"C",
			"Db",
			"D",
			"Eb",
			"E",
			"F",
			"Gb",
			"G",
			"Ab",
		}

		for i, key := range k {
			if i >= len(keys) {
				break
			}
			if slices.Contains(checklist, strings.Title(key)) {
				keys[i] = key
				continue
			}
			keys[i] = ""
		}
	}

	htmlString := RenderContent(app.reaper.Sections, keys, app.renderer, displayChords, displayLyrics)
	app.logger.Debug("rendering autoscroll", "chords", displayChords, "lyrics", displayLyrics)

	page := Page{
		Page: Meta{
			Title: "Lyrics will scroll automatically",
			Stats: []stat{
				{
					Title: "IP address of the server",
					Value: fmt.Sprintf("%s:%s", app.IP, app.Port),
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

func RenderContent(data []reaper.Section, keys []string, render *templates.Templates, withChords, withLyrics bool) string {
	html := bytes.NewBufferString("")

	for _, s := range data {
		content := bytes.NewBufferString("")

		render.RenderComponent(content, s.Name, "title")

		render.RenderComponent(content, s.Comments, "comments")

		if withChords && withLyrics {
			RenderChorded(s.Lines, content, render, keys[s.ItemID])
		} else if withLyrics {
			RenderChordless(s.Lines, content, render)
		} else if withChords {
			RenderChords(s.Lines, content, render, keys[s.ItemID])
		} else {
			RenderNotes(s.Lines, content, render)
		}

		fmt.Fprintf(html, `<div id="section-%d">%s</div>`, s.ID, content)
	}

	return html.String()
}

func RenderChorded(lines lyrics.Lines, buf *bytes.Buffer, render *templates.Templates, key string) *bytes.Buffer {
	for _, line := range lines {
		render.RenderComponent(
			buf,
			struct {
				ID   int
				Text interface{}
			}{
				ID:   line.ID,
				Text: template.HTML(line.WithChordsHTMLKey(key)),
			}, "line")
	}

	return buf
}

func RenderChordless(lines lyrics.Lines, buf *bytes.Buffer, render *templates.Templates) *bytes.Buffer {
	for _, line := range lines {
		render.RenderComponent(
			buf,
			struct {
				ID   int
				Text interface{}
			}{
				ID:   line.ID,
				Text: strings.ToUpper(line.Chordless()),
			}, "line")
	}

	return buf
}

func RenderChords(lines lyrics.Lines, buf *bytes.Buffer, render *templates.Templates, key string) *bytes.Buffer {
	for _, line := range lines {
		render.RenderComponent(
			buf,
			struct {
				ID   int
				Text interface{}
			}{
				ID:   line.ID,
				Text: template.HTML(line.Chords.HTMLKey(key)),
			}, "line")
	}

	return buf
}

func RenderNotes(lines lyrics.Lines, buf *bytes.Buffer, render *templates.Templates) *bytes.Buffer {
	for _, line := range lines {
		render.RenderComponent(
			buf,
			struct {
				ID   int
				Text interface{}
			}{
				ID:   line.ID,
				Text: "",
			}, "line")
	}

	return buf

}
