package main

import (
	"fmt"
	"log"
	"net/http"
)

type SectionType string
type Section struct {
	ID    int
	Name  string
	Text  []string
	Type  SectionType
	Image string
}

const (
	TypeSection SectionType = "section"
	TypeHero    SectionType = "hero"
)

func (app *app) autoscrollHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	file := "autoscroll.html"

	data := app.CreateMap()

	page := Page{
		Page: Meta{
			Title: "Lyrics will scroll automatically",
		},
		Data: data,
	}

	for _, item := range data {

		if item.Type != TypeHero {
			continue
		}

		page.Page.TableOfContent = append(page.Page.TableOfContent, TableOfContentEntry{
			Title: item.Name,
			Link:  fmt.Sprintf("#%d", item.ID),
		})

	}

	switch r.Method {
	case "POST":
		err = app.renderer.RenderComponent(w, page.Data, "autoscroll")
	default:
		err = app.renderer.RenderPage(w, page, file)

	}

	if err != nil {
		log.Println("failed to execute template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO: There is a known bug in CreateMap function in the section when I determent wether the
// region type is hero or not. If item splits on the section it will become a hero type even if it's
// not the beginning of a song.
// for example if I cut an paste\mix match regions on the fly in the reaper project during rehersal
// every time there is a cut section will get a hero type and thus i'll lose lyrics for that section
func (app *app) CreateMap() []Section {
	var region Section
	var results []Section
	var songMap string

	for _, m := range app.reaper.Markers {
		if !m.IsRegion {

			continue

		}

		region.ID = m.ID
		region.Name = m.Name
		region.Type = TypeSection

		if len(m.Name) > 0 {
			songMap += string(m.Name[0]) + " "
		}

		for _, item := range app.reaper.Items {

			if item.Contains(m.Start) {
				if m.Start == item.Position {
					region.Type = TypeHero
					region.Image = item.Image
				}
				section, err := item.GetSection(m.Name)

				if err != nil {
					region.Text = []string{""}
					continue
				} else {
					var text []string
					for _, line := range section.Text {

						if line != "" {
							text = append(text, line)
						}

					}
					region.Text = text
				}
			}
		}

		results = append(results, region)
	}

	return results
}
