package main

import (
	"fmt"
	"net/http"

	"github.com/butbkadrug/reaper-lyrics/pkg/osc"
)

func (app *app) eventsHandler(w http.ResponseWriter, r *http.Request) {

	//we storing connections in the map using port as a key
	id := r.RemoteAddr
	ch := app.SSEClients.AddClient(id)

	//setting headder required by SSE standarts
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	//defering cleanup of the clients
	// I have a feeling that this is where the iessue lays but I don't know
	defer func() {
		app.SSEClients.RemoveClient(id)
	}()

	//checking if connection supports flushing
	flusher, ok := w.(http.Flusher)

	if !ok {
		app.logger.Error("sse is not supported", "id", id)
		http.Error(w, "sse not supported", http.StatusInternalServerError)
		return
	}

	//main SSE loop
	for val := range *ch {
		if err := app.processMessage(id, val, w, flusher); err != nil {
			break
		}
	}
}

func (app *app) processMessage(id string, val interface{}, w http.ResponseWriter, flusher http.Flusher) error {
	switch val.(type) {
	case osc.PositionUpdate:
		_, err := fmt.Fprintf(w, "event: line\ndata:%d\n\n", app.reaper.Line.ID)
		if err != nil {
			app.logger.Error("events: Fprintf failed", "error", err.Error(), "id", id)
			return err
		}
		flusher.Flush()
	case osc.StopUpdate:
		_, err := fmt.Fprintf(w, "event: stop\ndata:%s\n\n", "playback stopped")
		if err != nil {
			app.logger.Error("events: Fprintf failed", "error", err.Error(), "id", id)
			return err
		}
		flusher.Flush()
	case osc.PlayUpdate:
		_, err := fmt.Fprintf(w, "event: play\ndata:%s\n\n", "playback started")
		if err != nil {
			app.logger.Error("events: Fprintf failed", "error", err.Error(), "id", id)
			return err
		}
		flusher.Flush()
	}

	return nil
}
