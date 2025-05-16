package storage

import (
	"zerodupe/internal/server/model"
)

// Storage defines the interface for storage operations
type Storage interface {
	// CheckFileExists checks if a file exists
	CheckFileExists(fileHash string) (bool, error)

	// CheckChunkExists checks if chunks exist
	CheckChunkExists(hashes []string) ([]string, []string, error)

	// SaveChunkMetadata saves chunk metadata
	SaveChunkMetadata(fileHash, chunkHash string, chunkOrder int) error

	// SaveChunkData saves chunk data
	SaveChunkData(chunkHash string, content []byte) (string, error)

	// GetFileMetadata gets file metadata
	GetFileMetadata(fileHash string) (*model.FileMetadata, error)

	// GetChunkData gets chunk data
	GetChunkData(chunkHash string) ([]byte, error)
}
