package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.autoscrollHandler)
	router.HandleFunc("/events", app.eventsHandler)
	router.HandleFunc("/stats", app.statsHandler)
	router.HandleFunc("GET /guitar", app.guitarHandler)
	router.HandleFunc("/guitar/", app.guitarHandler)
	router.HandleFunc("GET /drums", app.drumsHandler)
	router.HandleFunc("/drums/", app.drumsHandler)
	router.HandleFunc("/comments", app.commentHandler)
	router.HandleFunc("/comments/", app.commentHandler)
	router.HandleFunc("/htmx/set-keys", app.htmxSetKeysHandler)

	router.Handle("/static/", http.FileServer(http.FS(app.renderer.Static)))

	return app.loggerFunc(router)
}
