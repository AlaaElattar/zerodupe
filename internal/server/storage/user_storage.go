package storage

import "zerodupe/internal/server/model"

type UserStorage interface {
	CreateUser(user *model.User, plainPassword string) error
	LoginUser(username, password string) (*model.User, error)
}
