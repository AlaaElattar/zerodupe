package storage

// FileSystem defines the interface for storage operations
type FileSystem interface {
	// CheckFileExists checks if a file exists
	CheckFileExists(fileHash string) (bool, error)

	// CheckChunkExists checks if chunks exist
	CheckChunkExists(hashes []string) ([]string, []string, error)

	// SaveChunkData saves chunk data
	SaveChunkData(chunkHash string, content []byte) (string, error)

	// GetChunkData gets chunk data
	GetChunkData(chunkHash string) ([]byte, error)
}
