package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *app) eventsHandler(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "sse not supported", http.StatusInternalServerError)
		log.Println("failed to cast flusher")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
			case id := <-app.osc:

			log.Println("new message: ", id)

			n, err := fmt.Fprintf(w, "data: %s\n\n", id)

            if err != nil {
                log.Println("failed to write response: ", err)
                break
            }

			log.Printf("bytes writen %d. ERROR: %v\n", n, err)

			flusher.Flush()

		}
	}
}
