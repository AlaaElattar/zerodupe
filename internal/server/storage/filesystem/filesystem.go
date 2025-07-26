package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"zerodupe/pkg/hasher"

	"github.com/rs/zerolog/log"
)

// FilesystemStorage implements the Storage interface using the filesystem
type FilesystemStorage struct {
	storageDir string
}

// NewFilesystemStorage creates a new filesystem storage
func NewFilesystemStorage(storageDir string) (*FilesystemStorage, error) {
	// Create storage directories if they don't exist
	dirs := []string{
		storageDir,
		filepath.Join(storageDir, "blocks"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return &FilesystemStorage{
		storageDir: storageDir,
	}, nil
}

// CheckFileExists checks if a file exists in meta or blocks directory
func (fs *FilesystemStorage) CheckFileExists(fileHash string) (bool, error) {
	// check if it's a single chunk file
	blockPath := filepath.Join(fs.storageDir, "blocks", fileHash[:4], fileHash)
	if _, err := os.Stat(blockPath); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}

	return false, nil
}

// CheckChunkExists checks if chunks exist
func (fs *FilesystemStorage) CheckChunkExists(hashes []string) ([]string, []string, error) {
	var existingChunks []string
	var missingChunks []string

	for _, hash := range hashes {
		blockPath := filepath.Join(fs.storageDir, "blocks", hash[:4], hash)

		_, err := os.Stat(blockPath)
		if err != nil {
			if os.IsNotExist(err) {
				missingChunks = append(missingChunks, hash)
			} else {
				return nil, nil, err
			}
		} else {
			existingChunks = append(existingChunks, hash)
		}
	}

	return existingChunks, missingChunks, nil
}

// SaveChunkData saves chunk data
func (fs *FilesystemStorage) SaveChunkData(chunkHash string, content []byte) (string, error) {
	// Ensure directory exists
	blockDir := filepath.Join(fs.storageDir, "blocks", chunkHash[:4])
	if err := os.MkdirAll(blockDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create block directory: %w", err)
	}

	blockPath := filepath.Join(blockDir, chunkHash)

	// Check if chunk already exists
	if _, err := os.Stat(blockPath); err == nil {
		return chunkHash, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check if chunk exists: %w", err)
	}

	// Verify chunk hash
	isValid, calculatedHash := hasher.VerifyChunkHash(content, chunkHash)
	if !isValid {
		log.Warn().Msgf("Hash mismatch. Expected: %s, Got: %s", chunkHash, calculatedHash)
	}

	// Write chunk data
	if err := os.WriteFile(blockPath, content, 0644); err != nil {
		return calculatedHash, fmt.Errorf("failed to write chunk data: %w", err)
	}

	return calculatedHash, nil
}

// GetChunkData gets chunk data
func (fs *FilesystemStorage) GetChunkData(chunkHash string) ([]byte, error) {
	blockPath := filepath.Join(fs.storageDir, "blocks", chunkHash[:4], chunkHash)

	if _, err := os.Stat(blockPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("chunk not found: %w", err)
	}

	content, err := os.ReadFile(blockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk data: %w", err)
	}
	log.Debug().Msgf("Read chunk %s, size: %d bytes", chunkHash, len(content))
	return content, nil
}
