package sqlite

import (
	"zerodupe/internal/server/auth"
	"zerodupe/internal/server/model"

	"gorm.io/gorm"
)

// SqliteMock implements the UserStorage interface for testing
type SqliteMock struct {
	users map[string]*model.User
}

func NewSqliteStorageMock() (*SqliteMock) {
	return &SqliteMock{
		users: make(map[string]*model.User),
	}
}

// CreateUser creates a new user in the mock storage
func (s *SqliteMock) CreateUser(user *model.User, plainPassword string) error {
	existing, err := s.GetUserByUsername(user.Username)
	if err == nil && existing != nil {
		return gorm.ErrDuplicatedKey
	}

	salt, err := auth.GenerateSalt()
	if err != nil {
		return err
	}

	hashed := auth.HashPassword(plainPassword, salt)
	user.Password = hashed
	user.Salt = salt

	s.users[user.Username] = user
	return nil
}

// LoginUser logs in a user from the mock storage
func (s *SqliteMock) LoginUser(username, password string) (*model.User, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	hashedPassword := auth.HashPassword(password, user.Salt)
	if hashedPassword != user.Password {
		return nil, gorm.ErrRecordNotFound
	}

	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Salt:     user.Salt,
	}, nil
}

// GetUserByUsername gets a user by username from the mock storage
func(s *SqliteMock) GetUserByUsername(username string) (*model.User, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Salt:     user.Salt,
	}, nil

}
