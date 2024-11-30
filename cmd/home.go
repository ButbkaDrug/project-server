package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"lp/pkg/comments"

	"github.com/butbkadrug/reaper-lyrics/pkg/reaper"
)

type MapItem struct {
	//referance to the line on the page
	//this id should represent the id of the
	//line item attaches to
	LineID int
	//reference to item original position in the
	//list of items
	ItemID int
	//string that represents item's name. It will
	//be displayed on the page
	Name string
}

type Meta struct {
	Title       string
	Description string
	Items       []MapItem
	Stats       []stat
}

type PageConfig struct {
	CommentsType string
	CommentsBig  bool
	WithChords   bool
	WithLyrics   bool
	Renderer     Renderer
	Keys         []string
	Sections     []reaper.Section
}

type Page struct {
	Page   Meta
	Data   interface{}
	Config PageConfig
}

type Renderer interface {
	RenderComponent(io.Writer, interface{}, string) error
}

// 4e izmeriali li Ivanova Valiantsina. I kak pereslat' den'gi
type PageBuilder interface {
	WithLyrics() PageBuilder
	WithChords() PageBuilder
	//change parametr to enum
	WithBigComments() PageBuilder
	//probobly change param to enum
	SetComments(string) PageBuilder
	//probobly also should change to slice of enum values
	SetKeys([]string) PageBuilder
	SetRenderer(Renderer) PageBuilder
	SetTitle(string) PageBuilder
	SetDescription(string) PageBuilder
	SetItems([]MapItem) PageBuilder
	SetData([]reaper.Section) PageBuilder
	Build() (Page, error)
}

type StagePageBuilder struct {
	page Page
}

func (s *StagePageBuilder) WithLyrics() PageBuilder {
	s.page.Config.WithLyrics = true

	return s
}

func (s *StagePageBuilder) WithChords() PageBuilder {
	s.page.Config.WithChords = true
	return s
}

func (s *StagePageBuilder) WithBigComments() PageBuilder {
	s.page.Config.CommentsBig = true
	return s
}

func (s *StagePageBuilder) SetRenderer(r Renderer) PageBuilder {
	s.page.Config.Renderer = r

	return s
}

func (s *StagePageBuilder) SetTitle(t string) PageBuilder {
	s.page.Page.Title = t
	return s
}

func (s *StagePageBuilder) SetDescription(d string) PageBuilder {
	s.page.Page.Description = d
	return s
}

func (s *StagePageBuilder) SetItems(i []MapItem) PageBuilder {
	s.page.Page.Items = i
	return s
}

func (s *StagePageBuilder) SetComments(t string) PageBuilder {
	s.page.Config.CommentsType = t
	return s
}

func (s *StagePageBuilder) SetKeys(keys []string) PageBuilder {
	s.page.Config.Keys = keys

	return s
}

func (s *StagePageBuilder) SetData(sections []reaper.Section) PageBuilder {
	s.page.Config.Sections = sections
	return s
}

func (s *StagePageBuilder) Build() (Page, error) {

	html := bytes.NewBufferString("")
	render := s.page.Config.Renderer

	for _, sec := range s.page.Config.Sections {
		content := bytes.NewBufferString("")

		if err := render.RenderComponent(content, sec, "title"); err != nil {
			return Page{}, err
		}

		commentsStyle := "comments"

		if s.page.Config.CommentsBig {
			commentsStyle = "primer-comments"
		}

		comments := struct {
			ItemID    int
			SectionID int
			Comments  []comments.Comment
		}{
			ItemID:    sec.ItemID,
			SectionID: sec.ID,
			Comments:  GetCommentsOfType(sec.Comments, s.page.Config.CommentsType),
		}

		if err := render.RenderComponent(content, comments, commentsStyle); err != nil {

			return Page{}, err
		}

		if s.page.Config.WithChords && s.page.Config.WithLyrics {
			var key string
			if len(s.page.Config.Keys) > sec.ItemID {
				key = s.page.Config.Keys[sec.ItemID]
			}
			RenderChorded(sec.Lines, content, render, key)
		} else if s.page.Config.WithLyrics {
			RenderChordless(sec.Lines, content, render)
		} else if s.page.Config.WithChords {
			RenderChords(sec.Lines, content, render, s.page.Config.Keys[sec.ItemID])
		} else {
			RenderNotes(sec.Lines, content, render)
		}

		fmt.Fprintf(html, `<div id="section-%d">%s</div>`, sec.ID, content)
	}

	s.page.Data = struct {
		Content interface{}
	}{
		template.HTML(html.String()),
	}

	return s.page, nil
}

type PageMaster struct {
}

func (p *PageMaster) CreateSingersPage(b PageBuilder) (Page, error) {
	return b.SetTitle("WTEAM - Singer").
		SetDescription("No silly chords, no nonsense comments. Only good stuff all singers need!").
		WithLyrics().
		Build()

}

func (p *PageMaster) CreateChordedPage(b PageBuilder) (Page, error) {
	return b.SetTitle("WTEAM - Chords for nerds!").
		SetDescription("Page with chords. So, please use em!").
		WithChords().
		WithLyrics().
		Build()

}
