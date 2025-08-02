package storage

import "zerodupe/internal/server/model"

// UserStorage defines the interface for user storage operations
type UserStorage interface {
	// CreateUser creates a new user
	CreateUser(user *model.User) error

	// GetUserByUsername gets a user by username
	GetUserByUsername(username string) (*model.User, error)
}
