package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"zerodupe/internal/server/model"
	"zerodupe/pkg/hasher"
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
		filepath.Join(storageDir, "meta"),
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

// CheckFileExists checks if a file exists
func (fs *FilesystemStorage) CheckFileExists(fileHash string) (bool, error) {
	metaPath := filepath.Join(fs.storageDir, "meta", fileHash[:4], fileHash)
	if _, err := os.Stat(metaPath); err == nil {
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

// SaveChunkMetadata saves chunk metadata
func (fs *FilesystemStorage) SaveChunkMetadata(fileHash, chunkHash string, chunkOrder int) error {
	// Ensure directory exists
	metaDir := filepath.Join(fs.storageDir, "meta", fileHash[:4])
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	metaPath := filepath.Join(metaDir, fileHash)

	var metadata model.FileMetadata
	if _, err := os.Stat(metaPath); err == nil {
		content, err := os.ReadFile(metaPath)
		if err != nil {
			return fmt.Errorf("failed to read metadata file: %w", err)
		}
		if err := json.Unmarshal(content, &metadata); err != nil {
			return fmt.Errorf("failed to parse metadata: %w", err)
		}
	}

	// Check if chunk already exists in metadata
	for _, chunk := range metadata.Chunks {
		if chunk.ChunkOrder == chunkOrder && chunk.ChunkHash == chunkHash {
			return nil
		}
	}

	metadata.Chunks = append(metadata.Chunks, model.ChunkMetadataRequest{
		ChunkOrder: chunkOrder,
		ChunkHash:  chunkHash,
	})

	newContent, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
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
		fmt.Printf("WARNING: Hash mismatch. Expected: %s, Got: %s\n", chunkHash, calculatedHash)
	}

	// Write chunk data
	if err := os.WriteFile(blockPath, content, 0644); err != nil {
		return calculatedHash, fmt.Errorf("failed to write chunk data: %w", err)
	}

	return calculatedHash, nil
}

// GetFileMetadata gets file metadata
func (fs *FilesystemStorage) GetFileMetadata(fileHash string) (*model.FileMetadata, error) {
	metaPath := filepath.Join(fs.storageDir, "meta", fileHash[:4], fileHash)

	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file metadata not found: %w", err)
	}

	content, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata model.FileMetadata
	if err := json.Unmarshal(content, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &metadata, nil
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

	return content, nil
}
