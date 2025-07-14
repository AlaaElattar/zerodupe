package sqlite

import (
	"zerodupe/internal/server/model"

	"gorm.io/gorm"
)

// SqliteMock implements the UserStorage interface for testing
type SqliteMock struct {
	users map[string]*model.User
}

func NewSqliteStorageMock() *SqliteMock {
	return &SqliteMock{
		users: make(map[string]*model.User),
	}
}

// CreateUser creates a new user in the mock storage
func (s *SqliteMock) CreateUser(user *model.User) error {
	existing, err := s.GetUserByUsername(user.Username)
	if err == nil && existing != nil {
		return gorm.ErrDuplicatedKey
	}

	s.users[user.Username] = user
	return nil
}

// GetUserByUsername gets a user by username from the mock storage
func (s *SqliteMock) GetUserByUsername(username string) (*model.User, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil

}
