package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const shutdownPeriod = 15 * time.Second

var rootCmd = &cobra.Command{
	Use: "core",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to execute root command", "error", err)
	}
}

func getLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
