package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *app) eventsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.RemoteAddr
	ch := app.SSEClients.AddClient(id)

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
			fmt.Println("connection closed", r.RemoteAddr)
			return
		case <-app.updateCh:
			log.Println("file updated I'm ready to send an event!")
		default:
			val := <-*ch
			var err error
			switch val := val.(type) {
			case bool:
				_, err = fmt.Fprintf(w, "event: update\ndata: {\"update\":%v}\n\n", val)
				flusher.Flush()
			case string:
				_, err = fmt.Fprintf(w, "data: %s\n\n", val)
				flusher.Flush()
			case float32:
				item, _ := app.reaper.GetItem(float64(val))
				marker, _ := app.reaper.GetMarker(float64(val))
				section, _ := item.GetSection(marker.Name)

				numLines := len(section.Text)
				length := marker.End - marker.Start
				lineLengh := length / float64(numLines)
				lineNumber := int((float64(val) - marker.Start) / lineLengh)

				if lineNumber < 0 {
					lineNumber = 0
				}

				lastLine := fmt.Sprintf("%d-%d", marker.ID, lineNumber)

				app.reaper.LastLine = lastLine

				_, err = fmt.Fprintf(w, "event: cursor\ndata: {\"section\": %d, \"line\": %d}\n\n", marker.ID, lineNumber)
			}

			if err != nil {
				break
			}

			flusher.Flush()
		}

	}
}
