package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"sse"

	"github.com/butbkadrug/reaper-lyrics/pkg/osc"
	"github.com/butbkadrug/reaper-lyrics/pkg/projects"
	"github.com/butbkadrug/reaper-lyrics/pkg/reaper"
	"github.com/butbkadrug/reaper-lyrics/pkg/templates"
)

type app struct {
	server            *http.Server
	renderer          *templates.Templates
	oscClient         *osc.OSC
	defaultProjectDir string
	SSEClients        *sse.SSE
	startTime         time.Time
	reaper            *reaper.Project
	library           *projects.Library
	updateCh          chan bool
	logger            *slog.Logger
}

type config struct {
	// http server port to use
	port int
	// path to a reaper project
	projectPath string
	// should be true if it's not a host project that contains subprojects as songs but rather
	//song itself
	single bool
	// path to a setlist file
	setlistPath string
	// Where to look for a project folder/file if it is missing. In my case
	// when I share the same project directory between windows and linux
	// machine. So when I switch I want my thing to still find the lyrics
	// even if I didn't fix the project file yet.
	// It should save me a lot of worries.
	defaultProjectDir string
	// osc address to listen to
	oscIP string
	// osc port to listen to
	oscPort int
	//logger for the file
	logger *slog.Logger
}

func NewApp(cfg *config) *app {

	r := reaper.NewProject()

	if !cfg.single {
		r = reaper.NewProjectWithConfig(&reaper.ProjectConfig{
			Path:    cfg.projectPath,
			TrackID: 0,
			Logger:  cfg.logger,
		})

	} else {
		r = reaper.NewSongProjectWithConfig(&reaper.ProjectConfig{
			Path:    cfg.projectPath,
			TrackID: 0,
			Logger:  cfg.logger,
		})
	}

	clients := sse.NewSSE()

	renderer, err := templates.NewTemplate()

	if err != nil {
		log.Fatal("failed to start renderer. aborting...", err)
	}
	osc := osc.NewClient(&osc.OscCongin{
		IP:      cfg.oscIP,
		Port:    cfg.oscPort,
		Clients: clients,
		Logger:  cfg.logger,
	})
	lib := projects.NewLibrary(cfg.defaultProjectDir)

	app := &app{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", cfg.port),
		},
		renderer:          renderer,
		oscClient:         osc,
		reaper:            r,
		defaultProjectDir: cfg.defaultProjectDir,
		SSEClients:        clients,
		startTime:         time.Now(),
		library:           lib,
		logger:            cfg.logger,
	}

	cfg.logger.Debug("about to return new app")

	return app
}

func (app *app) getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}

	defer conn.Close()
	addr := conn.LocalAddr().(*net.UDPAddr)

	return fmt.Sprintf("%v", addr.IP)
}

func (app *app) serve() error {

	app.logger.Debug("app.go -> serve()")
	app.server.Handler = app.routes()

	go app.reaper.Watch(app.SSEClients)

	go app.library.ParseLibrary()

	go app.oscClient.Server.ListenAndServe()

	ip := app.getIP()

	msg := fmt.Sprintf("|%s|\n", strings.Repeat("-", 44))
	msg += fmt.Sprintf("|%s|\n", strings.Repeat(" ", 44))
	msg += fmt.Sprintf("| http server serves at: %s%s |\n", ip, app.server.Addr)
	msg += fmt.Sprintf("| osc client listens at: %s %s|\n", app.oscClient.Server.Addr, strings.Repeat(" ", 5))
	msg += fmt.Sprintf("|%s|\n", strings.Repeat(" ", 44))
	msg += fmt.Sprintf("|%s|\n", strings.Repeat("-", 44))

	fmt.Printf(msg)

	app.logger.Debug("about to start listening and serving")
	return app.server.ListenAndServe()
}
