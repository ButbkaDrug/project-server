package osc

import (
	"fmt"

	"github.com/butbkadrug/reaper-lyrics/pkg/sse"
	"github.com/hypebeast/go-osc/osc"
)

func NewClient(ip string, port int, clients *sse.SSE) *osc.Server{
    if ip == "" {
        ip = "127.0.0.1"
    }

    if port == 0 {
        port = 7400
    }

    addr := fmt.Sprintf("%s:%d", ip, port)
    d := osc.NewStandardDispatcher()


    d.AddMsgHandler("/time", func(msg *osc.Message){
        val, ok := msg.Arguments[0].(float32)

        if !ok {
            return
        }

        go clients.Broadcast(val)

    })

    return &osc.Server{
        Addr:       addr,
        Dispatcher: d,
    }
}
