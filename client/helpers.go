package client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// ChunkSize defines the size of each chunk in bytes (1MB)
const ChunkSize = 1 * 1024 * 1024

// FileChunk represents a single chunk of a file
type FileChunk struct {
	Data       []byte
	ChunkHash  string
	ChunkOrder int
}

// getFileChunks returns the SHA-256 hash of a file's content
func (client *Client) getFileChunks(filePath string) ([]FileChunk, string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("file does not exist: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, "", fmt.Errorf("failed to calculate file hash: %w", err)
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))

	if _, err := file.Seek(0, 0); err != nil {
		return nil, "", fmt.Errorf("failed to seek to beginning of file: %w", err)
	}

	var chunks []FileChunk
	buffer := make([]byte, ChunkSize)
	chunkOrder := 1

	// Splitting the file into chunks
	for {
		bytes, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, "", fmt.Errorf("failed to read file: %w", err)
		}
		if bytes == 0 {
			break
		}

		chunkData := make([]byte, bytes)
		copy(chunkData, buffer[:bytes])

		chunkHasher := sha256.New()
		chunkHasher.Write(chunkData)
		chunkHash := hex.EncodeToString(chunkHasher.Sum(nil))

		chunks = append(chunks, FileChunk{
			Data:       chunkData,
			ChunkHash:  chunkHash,
			ChunkOrder: chunkOrder,
		})

		chunkOrder++
		if err == io.EOF {
			break
		}
	}

	return chunks, fileHash, nil
}
