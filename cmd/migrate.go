package cmd

import (
	"context"
	"core/internal/database"
	"core/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate the database",
	Run:  migrate,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config, err := loadConfig()
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	l := getLogger()

	db, err := database.NewClient(ctx, config.DatabaseDSN)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.Migration00001Init,
	})

	if err = m.Migrate(); err != nil {
		l.Fatal("Migration failed", zap.Error(err))
	}

	l.Info("Migration finished successfully")
}
