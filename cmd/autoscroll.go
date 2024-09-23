// TODO: I want to move all the logic from template into here and have as little logic inside the telplate as possible
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"lp/pkg/comments"
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

func (app *app) ParseCapoArgument(s string) []int {
	capo := make([]int, len(app.reaper.Track.Items))

	if c := strings.Split(s, ","); len(c) > 0 {
		for i, entry := range c {

			if i >= len(capo) {
				break
			}

			if num, err := strconv.Atoi(entry); err != nil {
				capo[i] = 0
			} else {
				capo[i] = num
			}
		}

	}

	fmt.Println("hello from app.ParseCapoArgument(): ", capo)

	return capo
}

// parser song keys from a url params
func (app *app) ParseKeysArgument(s string) []string {
	keys := make([]string, len(app.reaper.Track.Items))

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

	if k := strings.Split(s, ","); len(k) > 0 {

		for i, key := range k {
			key = strings.ToUpper(key)
			if i >= len(keys) {
				break
			}
			if slices.Contains(checklist, key) {
				keys[i] = key
				continue
			}
			keys[i] = ""
		}
	}

	return keys
}

func GetCommentsOfType(c []comments.Comment, t string) []comments.Comment {

	t = strings.ToLower(t)

	fmt.Println("GetCommentsOfType", c, t)

	if t == "all" {
		return c
	}

	if t == "" {
		return comments.FilterComments(c, comments.GENERIC)
	}

	requestedTypes := strings.Split(t, " ")

	requestedComments := []comments.Comment{}

	for _, rt := range requestedTypes {

		trueType := comments.GetCommentType(rt)

		if trueType == comments.GENERIC && rt != "generic" {
			continue
		}

		requestedComments = append(requestedComments, comments.FilterComments(c, trueType)...)
	}

	if slices.Contains(requestedTypes, "only") {
		return requestedComments
	}

	//I want generic comments to be on top, that's why i reverse slices duren append
	return append(comments.FilterComments(c, comments.GENERIC), requestedComments...)

}

type RenderContentConfig struct {
	Data         []reaper.Section
	Render       *templates.Templates
	CommentsType string
	CommentsBig  bool
	WithChords   bool
	WithLyrics   bool
	Keys         []string
	Capo         []int
}

func RenderContent(cfg *RenderContentConfig) string {
	html := bytes.NewBufferString("")
	render := cfg.Render

	for _, s := range cfg.Data {
		content := bytes.NewBufferString("")

		render.RenderComponent(content, s.Name, "title")

		commentsStyle := "comments"

		if cfg.CommentsBig {
			commentsStyle = "primer-comments"
		}

		render.RenderComponent(content, GetCommentsOfType(s.Comments, cfg.CommentsType), commentsStyle)

		if cfg.WithChords && cfg.WithLyrics {
			RenderChorded(s.Lines, content, render, cfg.Keys[s.ItemID])
		} else if cfg.WithLyrics {
			RenderChordless(s.Lines, content, render)
		} else if cfg.WithChords {
			RenderChords(s.Lines, content, render, cfg.Keys[s.ItemID])
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
