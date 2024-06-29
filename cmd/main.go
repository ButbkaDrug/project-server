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
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {

	cfg, err := parseFlags()

	if err != nil {
		log.Fatalf("faild parsing arguments: %v", err)
	}

	logfile, _ := os.Create("project.log")
	defer logfile.Close()

	logger := slog.New(slog.NewTextHandler(logfile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg.logger = logger

	err = godotenv.Load()

	if err != nil {
		logger.Error("failed to load evirement variables", "error", err.Error())
	}

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

	flag.StringVar(
		&cfg.setlistPath,
		"setlist",
		cfg.setlistPath,
		"path to a rpp project to be parsed",
	)

	flag.Parse()

	return cfg, nil
}
