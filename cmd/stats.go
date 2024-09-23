package main

import (
	"fmt"
	"net/http"
	"time"
)

type stat struct {
	Title string
	Value any
	Desc  string
}

func (app *app) statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")

	file := "stats.html"

	data := []stat{

		{
			Title: "Location",
			Value: app.reaper.Path,
			Desc:  "Project file location",
		},

		{
			Title: "Modified",
			Value: time.Now().Sub(app.reaper.Stats[app.reaper.Path]),
			Desc:  "Since last changes were made",
		},

		{
			Title: "Connections",
			Value: len(app.SSEClients.GetClientsId()),
			Desc:  "clients listening for SSE",
		},

		{
			Title: "HTTP Addres",
			Value: fmt.Sprintf("%s%s", app.getIP(), app.server.Addr),
			Desc:  "IP addrres server serves at",
		},

		{
			Title: "OSC Addres",
			Value: app.oscClient.Server.Addr,
			Desc:  "address osc client listens on",
		},

		{
			Title: "Uptime",
			Value: time.Now().Sub(app.startTime),
			Desc:  "time since server loaded",
		},
	}

	page := Page{
		Page: Meta{
			Title: "Server Statistics",
		},
		Data: data,
	}

	err := app.renderer.RenderPage(w, page, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
