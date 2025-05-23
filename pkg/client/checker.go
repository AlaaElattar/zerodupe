package client

import (
	"fmt"
	"zerodupe/pkg/hasher"
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
	chunkHashes := extractChunkHashes(chunks)

	// Get missing chunks from server
	missingChunks, err := fc.api.GetMissingChunks(chunkHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to check chunks: %w", err)
	}

	// If all chunks exist, we can skip the upload
	if len(missingChunks) == 0 {
		fmt.Printf("All chunks already exist on server. Skipping upload.\n")
		return nil, nil
	}

	existingChunksMap := buildExistingChunksMap(chunkHashes, missingChunks)

	fmt.Printf("Existing chunks: %d, Missing chunks: %d\n", len(existingChunksMap), len(missingChunks))

	return existingChunksMap, nil
}

// extractChunkHashes extracts hashes from a slice of FileChunks
func extractChunkHashes(chunks []hasher.FileChunk) []string {
	hashes := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		hashes = append(hashes, chunk.ChunkHash)
	}
	return hashes
}

// buildExistingChunksMap creates a map of existing chunks by comparing all chunks with missing chunks
func buildExistingChunksMap(allHashes []string, missingHashes []string) map[string]bool {
	missingMap := make(map[string]bool)
	for _, hash := range missingHashes {
		missingMap[hash] = true
	}

	existingMap := make(map[string]bool)
	for _, hash := range allHashes {
		if !missingMap[hash] {
			existingMap[hash] = true
		}
	}

	return existingMap
}
