//NOTE: I think effective way to show the line is to know the offset since the
// beginning of the section. Then We don't have to know the exact order. And
// as long as we know which section we're in and what time. We can calculate
// which line is played.

package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"

	"github.com/joho/godotenv"
)

func main() {

	level := slog.LevelError
	if slices.Contains(os.Args, "DEBUG") {
		level = slog.LevelDebug
	}

	cfg, err := parseFlags()

	if err != nil {
		log.Fatalf("faild parsing arguments: %v", err)
	}

	logfile, _ := os.Create("project.log")
	defer logfile.Close()

	logger := slog.New(slog.NewTextHandler(logfile, &slog.HandlerOptions{
		Level: level,
	}))

	cfg.logger = logger

	err = godotenv.Load()

	if err != nil {
		logger.Error("failed to load evirement variables", "error", err.Error())
	}

	logger.Info("config for the project", "cfg", cfg)

	app := NewApp(&cfg)

	if err = app.serve(); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}

func parseFlags() (config, error) {
	var host string
	var oscPort int
	var cfg config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cfg, err
	}

	cfg.defaultProjectDir = filepath.Join(
		homeDir,
		"Music",
		"REAPER",
	)

	cfg.projectPath = filepath.Join(
		homeDir,
		"Music",
		"REAPER",
		"BaseProject",
		"BaseProject.rpp",
	)

	cfg.single = false

	cfg.setlistPath = filepath.Join(
		homeDir,
		"Music",
		"REAPER",
		"BaseProject",
		"setlist",
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
		"port of OSC host that we will listen on",
	)

	flag.IntVar(
		&cfg.port,
		"sp",
		1234,
		"port for http server to use for serving pages",
	)

	flag.BoolVar(
		&cfg.single,
		"single",
		false,
		"should be set if the project is a '.rpp' regular project(without subprojects)",
	)

	flag.Parse()

	for _, arg := range os.Args {
		//NOTE: If any argument has .rpp extension we will assume that it is a desired project path.
		if ext := path.Ext(arg); ext == ".rpp" || ext == ".RPP" {
			log.Printf("Project provided by arguments: %v", arg)
			cfg.projectPath = arg
		}
	}

	return cfg, nil
}
