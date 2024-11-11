package cmd

import (
	"context"
	"core/internal/database"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
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

	// Initialize database connection
	db, err := database.NewClient(ctx, config.DatabaseDSN)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}

	// Seed the database
	l.Info("seeding database")
	if err := seedDatabase(ctx, db); err != nil {
		l.Fatal("failed to seed database", zap.Error(err))
	}

	l.Info("done")
}

func seedDatabase(ctx context.Context, db *sqlx.DB) error {
	files, err := os.ReadDir("internal/database/seeds")
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".sql") {
			fileNames = append(fileNames, file.Name())
		}
	}
	sort.Strings(fileNames) // Sort the file names for deterministic order

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	for _, fileName := range fileNames {
		filePath := filepath.Join("seeds", fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			tx.Rollback()

			return fmt.Errorf("failed to read file %s: %w", fileName, err)
		}

		statements := strings.Split(string(content), ";")

		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}

			if _, err := tx.ExecContext(ctx, statement); err != nil {
				tx.Rollback()

				return fmt.Errorf("failed to execute statement from %s: %w", fileName, err)
			}
		}

		zap.L().Info("executed seed file", zap.String("file", fileName))
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
