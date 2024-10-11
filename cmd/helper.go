package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"lp/pkg/comments"
	"lp/pkg/lyrics"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/butbkadrug/reaper-lyrics/pkg/reaper"
	"github.com/butbkadrug/reaper-lyrics/pkg/templates"
)

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

	return capo
}

// parser song keys from a url params
func (app *app) ParseKeysArgument(v *url.Values) []string {
	//FIXME: i put arbitraty value for now, which i not very lickly to cross, but it shuld be equal to number of items that are in the project
	keys := make([]string, 10)

	s := v.Get("keys")

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
			key = strings.Title(key)
			if i >= len(k) {
				break
			}
			if slices.Contains(checklist, key) {
				keys[i] = key
				continue
			}
			keys[i] = ""
		}
	}

	for i := range 10 {
		key := fmt.Sprintf("key%d", i)

		if !v.Has(key) {
			continue
		}

		val := strings.Title(v.Get(key))

		if !slices.Contains(checklist, val) {
			continue
		}

		if i < 1 {
			continue
		}

		keys[i-1] = val
	}

	return keys
}

func GetCommentsOfType(c []comments.Comment, t string) []comments.Comment {

	t = strings.ToLower(t)

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
	logger       *slog.Logger
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

	if cfg == nil {
		panic("RenderContent invoket with nil config!")
	}

	html := bytes.NewBufferString("")
	render := cfg.Render

	for _, s := range cfg.Data {
		content := bytes.NewBufferString("")

		if err := render.RenderComponent(content, s, "title"); err != nil {
			cfg.logger.Error("cmd/helper.go RenderContent() failed", "error", err.Error(), "config", cfg)
		}

		commentsStyle := "comments"

		if cfg.CommentsBig {
			commentsStyle = "primer-comments"
		}

		comments := struct {
			ItemID    int
			SectionID int
			Comments  []comments.Comment
		}{
			ItemID:    s.ItemID,
			SectionID: s.ID,
			Comments:  GetCommentsOfType(s.Comments, cfg.CommentsType),
		}

		if err := render.RenderComponent(content, comments, commentsStyle); err != nil {

			cfg.logger.Error("cmd/helper.go -> renderComponent() *comments failed", "error", err.Error())
		}

		if cfg.WithChords && cfg.WithLyrics {
			var key string
			if len(cfg.Keys) > s.ItemID {
				key = cfg.Keys[s.ItemID]
			}
			RenderChorded(s.Lines, content, render, key)
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
