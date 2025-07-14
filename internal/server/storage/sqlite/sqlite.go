package sqlite

import (
	"zerodupe/internal/server/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SqliteStorage implements the UserStorage interface using SQLite
type Sqlite struct {
	db *gorm.DB
}

// NewSqliteStorage connects to the database file
func NewSqliteStorage(file string) (*Sqlite, error) {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	return &Sqlite{db: db}, nil
}

// Close closes the database connection
func (s *Sqlite) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateUser creates a new user
func (s *Sqlite) CreateUser(user *model.User) error {
	return s.db.Create(user).Error

}

// GetUserByUsername gets a user by username
func (s *Sqlite) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
