package cmd

import (
	"core/internal/database"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate the database to the newest migration",
	Run:  migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	config, err := loadConfig()
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	l := getLogger()

	l.Info("migrating database")

	db, err := database.NewClient(ctx, config.DatabaseDSN)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}

	if err := database.Migrate(ctx, db); err != nil {
		l.Fatal("failed to migrate database", zap.Error(err))
	}

	l.Info("done")
}
