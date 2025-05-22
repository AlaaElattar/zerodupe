package client

import (
	"fmt"
	"zerodupe/packages/hasher"
)

// FileChecker handles checking if files and chunks exist on the server
type FileChecker struct {
	api API
}

// NewFileChecker creates a new file checker
func NewFileChecker(api API) *FileChecker {
	return &FileChecker{
		api: api,
	}
}

// CheckFileExists checks if a file exists on the server
func (fc *FileChecker) CheckFileExists(fileHash string) (bool, error) {
	return fc.api.CheckFileExists(fileHash)
}

// IdentifyExistingChunks checks which chunks already exist on the server
func (fc *FileChecker) IdentifyExistingChunks(chunks []hasher.FileChunk) (map[string]bool, error) {
	// Extract chunk hashes
	chunkHashes := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		chunkHashes = append(chunkHashes, chunk.ChunkHash)
	}

	// Check which chunks exist on server
	existingChunks, missingChunks, err := fc.api.CheckChunksExists(chunkHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to check chunks: %w", err)
	}

	fmt.Printf("Existing chunks: %d, Missing chunks: %d\n", len(existingChunks), len(missingChunks))

	// If all chunks exist, we can skip the upload
	if len(missingChunks) == 0 {
		fmt.Printf("All chunks already exist on server. Skipping upload.\n")
		return nil, nil
	}

	// Create a map for quick lookup
	existingChunksMap := make(map[string]bool)
	for _, hash := range existingChunks {
		existingChunksMap[hash] = true
	}

	return existingChunksMap, nil
}
