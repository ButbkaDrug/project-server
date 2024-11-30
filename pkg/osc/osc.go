package osc

import (
	"errors"
	"fmt"

	"github.com/hypebeast/go-osc/osc"
)

const defaultPort = 7400

type PositionUpdate float64
type StopUpdate float32
type PlayUpdate float32

type Observer interface {
	Update(any)
}

type OSC struct {
	Server  *osc.Server
	clients []Observer
}

type options struct {
	port int
}

type Option func(*options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port must be a positive number!")
		}
		options.port = port
		return nil
	}
}

func NewClient(ip string, opts ...Option) (*OSC, error) {
	var options options

	for _, opt := range opts {
		err := opt(&options)

		if err != nil {
			return nil, err
		}
	}

	if ip == "" || ip == "localhost" {
		ip = "127.0.0.1"
	}

	var port int

	if options.port == 0 {
		port = defaultPort
	} else {
		port = options.port
	}

	addr := fmt.Sprintf("%s:%d", ip, port)

	server := &OSC{
		Server: &osc.Server{
			Addr: addr,
		},
	}

	d := osc.NewStandardDispatcher()

	err := d.AddMsgHandler("/time", func(msg *osc.Message) {
		val, ok := msg.Arguments[0].(float32)

		if !ok {
			return
		}

		server.notify(PositionUpdate(val))

	})

	if err != nil {
		return nil, err
	}

	err = d.AddMsgHandler("/play", func(msg *osc.Message) {
		val, ok := msg.Arguments[0].(float32)

		if !ok {
			return
		}

		if val == 0 {
			return
		}

		server.notify(PlayUpdate(val))

	})

	if err != nil {
		return nil, err
	}

	err = d.AddMsgHandler("/stop", func(msg *osc.Message) {
		val, ok := msg.Arguments[0].(float32)

		if !ok {
			return
		}

		if val == 0 {
			return
		}

		server.notify(StopUpdate(val))

	})

	if err != nil {
		fmt.Println("OSC ERROR:", err)
	}

	server.Server.Dispatcher = d

	return server, nil
}

func (osc *OSC) notify(v any) {
	for _, c := range osc.clients {
		c.Update(v)
	}
}

func (osc *OSC) Register(client Observer) {
	osc.clients = append(osc.clients, client)
}
