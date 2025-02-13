package database

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

//go:embed seeds/*.sql
var seeds embed.FS

func NewClient(ctx context.Context, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(ctx context.Context, db *sqlx.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}

func Seed(ctx context.Context, db *sqlx.DB) error {
	files, err := seeds.ReadDir("seeds")
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	var seedFiles []string
	for _, entry := range files {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			seedFiles = append(seedFiles, entry.Name())
		}
	}
	sort.Strings(seedFiles)

	for _, fileName := range seedFiles {
		content, err := os.ReadFile(filepath.Join("seeds", fileName))
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %w", fileName, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute seed file %s: %w", fileName, err)
		}
	}

	return nil
}
