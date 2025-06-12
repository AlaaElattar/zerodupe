package sqlite

import (
	"sync"
	"zerodupe/internal/server/auth"
	"zerodupe/internal/server/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SqliteStorage implements the UserStorage interface using SQLite
type Sqlite struct {
	db    *gorm.DB
	mutex sync.Mutex
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
func (s *Sqlite) CreateUser(user *model.User, plainPassword string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	salt, err := auth.GenerateSalt()
	if err != nil {
		return err
	}

	hashed := auth.HashPassword(plainPassword, salt)
	user.Password = hashed
	user.Salt = salt

	return s.db.Create(user).Error

}

// LoginUser logs in a user
func (s *Sqlite) LoginUser(username, password string) (*model.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	hashedPassword := auth.HashPassword(password, user.Salt)
	if hashedPassword != user.Password {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil

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
