package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockClient(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)

	columns := []string{"version"}
	mockSQL.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mockSQL.NewRows(columns).FromCSVString("1"),
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return db, mockSQL
}
