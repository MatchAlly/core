package database

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewClient(ctx context.Context, dsn string) (*gorm.DB, error) {
	dialect := mysql.Open(dsn)

	gormConfig := &gorm.Config{}
	gormConfig.Logger = logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
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
