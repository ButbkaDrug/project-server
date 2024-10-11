package osc

import (
	"fmt"
	"log/slog"

	"github.com/hypebeast/go-osc/osc"
	"sse"
)

type OSC struct {
	Server  *osc.Server
	Logger  *slog.Logger
	Clients *sse.SSE
}

type OscCongin struct {
	IP      string
	Port    int
	Logger  *slog.Logger
	Clients *sse.SSE
}

func NewClient(cfg *OscCongin) *OSC {
	if cfg.IP == "" {
		cfg.IP = "127.0.0.1"
	}

	if cfg.Port == 0 {
		cfg.Port = 7400
	}

	addr := fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	d := osc.NewStandardDispatcher()

	err := d.AddMsgHandler("/time", func(msg *osc.Message) {
		val, ok := msg.Arguments[0].(float32)

		cfg.Logger.Debug("OSC DISPATCHER: got /time message", "value", val)

		if !ok {
			return
		}

		go cfg.Clients.Broadcast(val)

	})

	if err != nil {
		cfg.Logger.Error("pkg/osc/osc.go -> failed to addMsgHandler(/time)", "error", err.Error())
	}

	return &OSC{
		Server: &osc.Server{
			Addr:       addr,
			Dispatcher: d,
		},
		Logger:  cfg.Logger,
		Clients: cfg.Clients,
	}
}
