package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

func (app *app) mapHandler(w http.ResponseWriter, r *http.Request) {
	f := path.Join("pkg", "templates", "songmap.gotmpl")

	t, err := template.ParseFS(templateMap, f)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type section struct {
		ID   int
		Name string
		Text []string
	}

	var region section
	var results []section

	for _, m := range app.reaper.Markers {
		region.ID = m.ID
		region.Name = m.Name

		for _, entry := range app.reaper.Items {
			if entry.Contains(m.Start) {
				log.Println(entry.Name, entry.Position, entry.Length)

				s, err := entry.GetSection(m.Name)

				if err != nil {
					region.Text = []string{}
					continue
				}

				region.Text = s.Text
			}
		}

		results = append(results, region)

	}

	err = t.Execute(w, results)

	if err != nil {
		log.Println("failed to execute template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
