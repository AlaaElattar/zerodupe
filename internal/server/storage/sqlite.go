package storage

import (
	"gorm.io/driver/sqlite"
)

// NewSqliteStorage connects to the database file
func NewSqliteStorage(file string) (DB, error) {
	dialector := sqlite.Open(file)
	return NewGormStorage(dialector)
}
