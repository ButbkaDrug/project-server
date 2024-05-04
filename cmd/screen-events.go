package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func(app *app) screeneventsHandler(w http.ResponseWriter, r *http.Request) {
    id := r.RemoteAddr
    ch := app.screen.Clients.AddClient(id)
    defer app.screen.Clients.RemoveClient(id)

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)

    if ! ok {

        http.Error(w, "see not supported, aborting...", http.StatusInternalServerError)
        return

    }

    for {

        select {
		case <-r.Context().Done():
			fmt.Println("connection closed", r.RemoteAddr)
            return
		default:
			val := <-*ch

            v, ok := val.([]string)

            if !ok {
                log.Println("error: not slice of string")
                return
            }

            _, err := fmt.Fprintf(w, "data: <h1>%s</h1>\n\n", strings.Join(v, "</h1><h1>"))

			if err != nil {
				log.Println("failed to write response: ", err)
				break
			}

			flusher.Flush()
        }
    }

}
