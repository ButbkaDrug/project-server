package main

import (
	"fmt"
	"net/http"
)

func (app *app) eventsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.RemoteAddr
	ch := app.SSEClients.AddClient(id)
	app.logger.Debug("new listener added", "id", id)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	defer app.SSEClients.RemoveClient(id)
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "sse not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			app.logger.Debug("connection closed", "address", r.RemoteAddr)
			return
		case <-app.updateCh:
			app.logger.Info("file updated I'm ready to send an event!")
		default:
			val := <-*ch
			app.logger.Debug("EVENTS: defauilt case", "value", val)
			var err error
			switch val := val.(type) {
			case bool:
				_, err = fmt.Fprintf(w, "event: update\ndata: {\"update\":%v}\n\n", val)
				flusher.Flush()
			case string:
				_, err = fmt.Fprintf(w, "data: %s\n\n", val)
				flusher.Flush()
			case float32:
				//NOTE: I want to confirm that client received line update. And while line doesn't
				// change do not send more updates. Then start sending updates again until confirmation is received

				cursor := float64(val)

				app.logger.Debug("got cursor", "c", cursor)

				app.reaper.Update(cursor)

				app.logger.Debug("updates finished. dispatching event", "lineID", app.reaper.Line.ID)
				_, err = fmt.Fprintf(w, "event: line\ndata:%d\n\n", app.reaper.Line.ID)
				flusher.Flush()

				if err != nil {
					app.logger.Error("dispatching failed", "error", err.Error())
				}

				if app.reaper.Item.Length == 0 {
					sendCountdownEvent(w, app.reaper.PeekItem.Name, app.reaper.PeekItem.Position-app.reaper.Cursor)
				}

			}

			if err != nil {
				break
			}

		}

	}
}

func sendAlertEvent(w http.ResponseWriter, msg string) error {

	_, err := fmt.Fprintf(w, "event: alert\ndata:{\"message\":\"%s\"}\n\n", msg)
	w.(http.Flusher).Flush()

	return err
}

func sendCountdownEvent(w http.ResponseWriter, name string, time float64) error {

	sec := int(time) % 60
	min := int(time) / 60

	if len(name) > 20 {
		name = name[:20] + "..."
	}

	_, err := fmt.Fprintf(w, "event: cowntdown\ndata:{\"name\":\"%s\", \"min\":%v, \"sec\":%v}\n\n",
		name,
		min,
		sec,
	)
	w.(http.Flusher).Flush()

	return err
}
