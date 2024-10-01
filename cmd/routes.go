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
	router.HandleFunc("/comments", app.commentHandler)
	router.HandleFunc("/comments/", app.commentHandler)

	router.Handle("/static/", http.FileServer(http.FS(app.renderer.Static)))

	return router
}
