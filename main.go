// TODO: Add watching functionality for main project as well as for every
// sub item. So that if any of them change we rescan the whole thig and
// stay up to date with project changes

// TODO: Get rig of temp parser? And maybe make a legit one...

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hypebeast/go-osc/osc"
)

func main() {

	var host string
	var oscPort int

	var cfg config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("failed to read filesystem: ", err)
	}

	cfg.projectPath = filepath.Join(
		homeDir,
		"Music",
		"REAPER",
		"BaseProject",
		"BaseProject.rpp",
	)

	flag.StringVar(
		&host,
		"host",
		"localhost",
		"ip address of OSC host",
	)

	flag.IntVar(
		&oscPort,
		"op",
		7400,
		"port of OSC host",
	)

	flag.IntVar(
		&cfg.port,
		"sp",
		1234,
		"port for http server to use",
	)

	flag.StringVar(
		&cfg.projectPath,
		"path",
		cfg.projectPath,
		"path to a rpp project to be parsed",
	)

	flag.Parse()

	log.Printf("Starging server:\nSERVER-PORT: %d\nOSC-HOST: %s\nOSC-PORT: %d\nProject: %s\n",
		cfg.port,
		host,
		oscPort,
		cfg.projectPath,
	)

	app := NewApp(&cfg)

	go func() {
		addr := fmt.Sprintf("%s:%d", host, oscPort)
		d := osc.NewStandardDispatcher()
		d.AddMsgHandler("/lastregion/number/str", func(msg *osc.Message) {
			if val, ok := msg.Arguments[0].(string); !ok || val == "" {
				return
			}

			go func() {
				val := msg.Arguments[0].(string)
				app.osc <- val
			}()

		})

		server := &osc.Server{
			Addr:       addr,
			Dispatcher: d,
		}
		server.ListenAndServe()
	}()

	err = app.serve()

	if err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
