package database

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN string `mapstructure:"dsn" validate:"required"`
}

func NewClient(ctx context.Context, config Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}
	gormConfig.Logger = logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(postgres.Open(config.DSN), gormConfig)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to connect to database, dsn: %s", config.DSN))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get database connection")
	}

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return db, nil
}
