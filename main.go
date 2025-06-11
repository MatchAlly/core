package main

import (
	"core/cmd"
	"flag"
	"log/slog"
	"os"
	"time"
)

const shutdownPeriod = 15 * time.Second

func main() {
	l := getLogger()

	apiCmd := flag.NewFlagSet("api", flag.ExitOnError)

	if len(os.Args) < 2 {
		slog.Error("Expected 'api' command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "api":
		apiCmd.Parse(os.Args[2:])
		cmd.StartAPIserver(l)
	default:
		slog.Error("Unknown command", "command", os.Args[1])
		os.Exit(1)
	}
}

func getLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
