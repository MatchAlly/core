package cmd

import (
	"core/internal/database"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS

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

	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		l.Fatal("failed to set dialect", zap.Error(err))
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		l.Fatal("failed to migrate database", zap.Error(err))
	}

	l.Info("done")
}
