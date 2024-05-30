package database

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewClient(ctx context.Context, dsn string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}
	gormConfig.Logger = logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to connect to database, dsn: %s", dsn))
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
