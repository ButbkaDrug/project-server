package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	client "github.com/butbkadrug/reaper-lyrics/pkg/osc"
	"github.com/butbkadrug/reaper-lyrics/pkg/parser"
	"github.com/butbkadrug/reaper-lyrics/pkg/projects"
	"github.com/butbkadrug/reaper-lyrics/pkg/screen"
	"github.com/butbkadrug/reaper-lyrics/pkg/setlist"
	"github.com/butbkadrug/reaper-lyrics/pkg/sse"
	"github.com/butbkadrug/reaper-lyrics/pkg/templates"
	"github.com/hypebeast/go-osc/osc"
)

type app struct {
	server            *http.Server
	renderer          *templates.Templates
	oscClient         *osc.Server
	defaultProjectDir string
	cursor            float64
	mu                *sync.Mutex
	SSEClients        *sse.SSE
	startTime         time.Time
	reaper            *parser.Project
	setlist           *setlist.Setlist
	library           *projects.Library
	screen            *screen.Screen
	updateCh          chan bool
}

type config struct {
	// http server port to use
	port int
	// path to a reaper project
	projectPath string
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
}

func NewApp(cfg *config) *app {

	r := parser.NewProject()
	r.Path = cfg.projectPath
	s := setlist.NewSetlist(cfg.setlistPath)
	clients := sse.NewSSE()

	app := &app{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", cfg.port),
		},
		renderer:          templates.NewTemplate(),
		oscClient:         client.NewClient(cfg.oscIP, cfg.oscPort, clients),
		reaper:            r,
		defaultProjectDir: cfg.defaultProjectDir,
		cursor:            0,
		mu:                &sync.Mutex{},
		SSEClients:        clients,
		startTime:         time.Now(),
		setlist:           s,
		library:           projects.NewLibrary(cfg.defaultProjectDir),
		screen:            screen.NewScreen(),
	}

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

	app.server.Handler = app.routes()
	err := app.reaper.ParseReaperFile()

	if err != nil {
		return fmt.Errorf("failed to parse project %s: %w", app.reaper.Path, err)
	}

	err = app.setlist.Read()

	if err != nil {
		return fmt.Errorf("failed to parse setlist %s: %w", app.setlist.Path, err)
	}

	go app.reaper.Watch(app.SSEClients)

	go app.library.ParseLibrary()

	go app.oscClient.ListenAndServe()

	ip := app.getIP()

	msg := fmt.Sprintf("|%s|\n", strings.Repeat("-", 44))
	msg += fmt.Sprintf("|%s|\n", strings.Repeat(" ", 44))
	msg += fmt.Sprintf("| http server serves at: %s%s |\n", ip, app.server.Addr)
	msg += fmt.Sprintf("| osc client listens at: %s %s|\n", app.oscClient.Addr, strings.Repeat(" ", 5))
	msg += fmt.Sprintf("|%s|\n", strings.Repeat(" ", 44))
	msg += fmt.Sprintf("|%s|\n", strings.Repeat("-", 44))

	fmt.Printf(msg)

	return app.server.ListenAndServe()
}
