package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.homeHandler)
	router.HandleFunc("/map", app.mapHandler)
	router.HandleFunc("/events", app.eventsHandler)
	router.Handle("/static/", http.FileServer(http.FS(static)))

	return router
}
