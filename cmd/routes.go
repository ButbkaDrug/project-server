package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.autoscrollHandler)
	router.HandleFunc("/events", app.eventsHandler)
	router.HandleFunc("/no-scroll", app.noscrollHandler)
	router.HandleFunc("/stats", app.statsHandler)
	router.HandleFunc("/songs", app.songsHandler)
	router.HandleFunc("/songs/{name}", app.viewsongHandler)

	router.HandleFunc("/setlist", app.setlistHandler)
	router.HandleFunc("POST /update-setlist", app.updateSetlist)
	router.HandleFunc("POST /add-to-setlist", app.addtosetlistHandler)

	router.HandleFunc("GET /screen-events", app.screeneventsHandler)
	router.HandleFunc("GET /screen", app.screenHandler)
	router.HandleFunc("POST /update-screen", app.updatescreenHandler)

	router.HandleFunc("GET /screen-master/{id}", app.screenmasterHandler)
	router.HandleFunc("GET /screen-master", app.screenmasterHandler)
	router.HandleFunc("/project-search", app.projectsearchHandler)

	router.HandleFunc("POST /get-project", app.getProjectHandler)

	router.Handle("/static/", http.FileServer(http.FS(app.renderer.Static)))

	return app.recoverPanic(router)
}
