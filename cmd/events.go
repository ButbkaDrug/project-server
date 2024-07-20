package main

import (
	"fmt"
	"net/http"
)

func (app *app) eventsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.RemoteAddr
	ch := app.SSEClients.AddClient(id)
	app.logger.Debug("new listener added", "id", id)

	params := r.URL.Query()

	app.logger.Debug("events handler. client subscribed", "params", params)

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

	// go func(w http.ResponseWriter, r *http.Request) {
	// 	select {
	// 	case <-time.After(200 * time.Millisecond):
	// 		sendLineUpdate(w, app.reaper.Line.ID)
	// 	}
	//
	// }(w, r)

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
				//TODO: I can assign id to every client when it connects and send it down with a message. and then when i respond with this id. i can stop sending the update client alread have.

				//NOTE: I want to confirm that client received line update. And while line doesn't
				// change do not send more updates. Then start sending updates again until
				// confirmation is receive. Stuck on that one though because, when client registers
				// for updates it has certain ip. but when it sends confirmation, it has different
				// ip which is changing all the time. which makes it difficult.

				//I shuld probobly issign id to each client and respond with id and confirmation
				//then changing ip will not metter.

				//one thing though that the whole strategy worked relly well for htmx. and now I don't relly use htmx in this part of application... but since i actually don't have optimize yet. we'll live with that.

				cursor := float64(val)

				app.logger.Debug("got cursor", "c", cursor)

				//TODO: Update shuld return a map\slice of updates that happened and base on those updates i can dispatch nessesery events
				app.reaper.Update(cursor)

				app.logger.Debug("dispatching lineID event", "lineID", app.reaper.Line.ID)
				_, err := fmt.Fprintf(w, "event: line\ndata:%d\n\n", app.reaper.Line.ID)
				flusher.Flush()

				if err != nil {
					app.logger.Error("dispatching failed", "error", err.Error())
				}

				// var data bytes.Buffer
				//
				// app.renderer.RenderComponent(&data, app.reaper.GetMapForCurrentItem(), "map")
				//
				// d := strings.ReplaceAll(data.String(), "\n", "")
				//
				// fmt.Fprintf(w, "event:songmap\ndata:%s\n\n", d)
				// flusher.Flush()
				//
				// fmt.Fprintf(w, "event: map\ndata:%d", app.reaper.Section.ID)
				// flusher.Flush()
				//
				// if app.reaper.Item.Length == 0 {
				// 	sendCountdownEvent(w, app.reaper.PeekItem.Name, app.reaper.PeekItem.Position-app.reaper.Cursor)
				// }

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
