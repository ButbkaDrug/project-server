package main

import (
	"embed"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/butbkadrug/reaper-lyrics/pkg/parser"
)

//go:embed pkg/templates/main.gotmpl
var templateMain embed.FS

//go:embed pkg/templates/songmap.gotmpl
var templateMap embed.FS

//go:embed static
var static embed.FS

type app struct {
	server *http.Server
	reaper *parser.Project
	osc    chan string
	cursor float64
	mu     *sync.Mutex
}

type config struct {
	// http server port to use
	port int
	// path to a reaper project
	projectPath string
	// http server ReadTimeout optional
	ReadTimeout time.Duration
	// http server WriteTimeout optional
	WriteTimeout time.Duration
	// http server IdleTimeout optional
	IdleTimeout time.Duration
}

func NewApp(cfg *config) *app {

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 5 * time.Second
	}


	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 30 * time.Second
	}

    r := parser.NewProject()
    r.Path = cfg.projectPath

	app := &app{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.port),
			ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: 0,
			IdleTimeout:  cfg.IdleTimeout,
		},
		reaper: r,
		osc:    make(chan string, 0),
		cursor: 0,
		mu:     &sync.Mutex{},
	}

	return app
}

func (app *app) serve() error {

	app.server.Handler = app.routes()
	err := app.reaper.ParseReaperFile()

	if err != nil {
		return fmt.Errorf("failed to parse project %s: %w", app.reaper.Path, err)
	}

	return app.server.ListenAndServe()
}
