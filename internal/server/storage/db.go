package storage

import "zerodupe/internal/server/model"

// DB defines the interface for db storage operations
type DB interface {
	// CreateUser creates a new user
	CreateUser(user *model.User) error

	// GetUserByUsername gets a user by username
	GetUserByUsername(username string) (*model.User, error)

	// SaveChunkMetadata saves chunk metadata
	SaveChunkMetadata(fileHash, chunkHash string, chunkOrder int) error

	// GetFileMetadata gets file metadata
	GetFileMetadata(fileHash string) (*model.FileMetadata, error)

	// CheckChunkExists checks if chunks exist in metadata
	CheckFileExists(fileHash string) (bool, error)
}
