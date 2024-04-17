package database

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AnyTime is used to match any time.Time value in tests
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

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
