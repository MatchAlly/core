package cmd

import (
	"core/internal/database"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var seedCmd = &cobra.Command{
	Use:  "seed",
	Long: "Seed the database",
	Run:  seed,
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func seed(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	config, err := loadConfig()
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	l := getLogger()

	db, err := database.NewClient(ctx, config.DatabaseDSN)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}

	if err := database.Seed(ctx, db); err != nil {
		l.Fatal("failed to seed database", zap.Error(err))
	}

	l.Info("done")
}
