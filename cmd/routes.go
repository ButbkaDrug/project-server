package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.autoscrollHandler)
	router.HandleFunc("/events", app.eventsHandler)
	router.HandleFunc("/stats", app.statsHandler)
	router.HandleFunc("/guitar", app.guitarHandler)
	router.HandleFunc("/guitar/", app.guitarHandler)
	router.HandleFunc("/drums", app.drumsHandler)
	router.HandleFunc("/drums/", app.drumsHandler)
	//TODO: next 2 should be refactored to benefit from lyrics parser
	router.HandleFunc("/songs", app.songsHandler)
	router.HandleFunc("/songs/{name}", app.viewsongHandler)

	router.Handle("/static/", http.FileServer(http.FS(app.renderer.Static)))

	return app.recoverPanic(router)
}
