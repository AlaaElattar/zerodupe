package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// ChunkSizeBytes defines the size of each chunk in bytes (1MB)
const ChunkSizeBytes = 1 * 1024 * 1024

// FileChunk represents a single chunk of a file
type FileChunk struct {
	Data       []byte
	ChunkHash  string
	ChunkOrder int
}

// CalculateChunkHash computes the SHA-256 hash of a byte slice
func CalculateChunkHash(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// SplitDataIntoChunks splits a byte slice into chunks and returns them along with the file hash
func SplitDataIntoChunks(data []byte) ([]FileChunk, string, error) {
	fileHasher := sha256.New()
	var chunks []FileChunk

	for i, order := 0, 1; i < len(data); i, order = i+ChunkSizeBytes, order+1 {
		end := i + ChunkSizeBytes
		if end > len(data) {
			end = len(data)
		}

		chunkData := data[i:end]
		chunkHash := CalculateChunkHash(chunkData)
		fileHasher.Write([]byte(chunkHash))

		chunks = append(chunks, FileChunk{
			Data:       chunkData,
			ChunkHash:  chunkHash,
			ChunkOrder: order,
		})
	}

	fileHash := hex.EncodeToString(fileHasher.Sum(nil))
	return chunks, fileHash, nil
}

// VerifyChunkHash verifies that a chunk's data matches its expected hash
func VerifyChunkHash(data []byte, expectedHash string) (bool, string) {
	actualHash := CalculateChunkHash(data)
	return actualHash == expectedHash, actualHash
}

func CombineChunksIntoFile(chunks [][]byte, outputDir string, fileName string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the output file
	outputPath := filepath.Join(outputDir, fileName)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	for i, chunk := range chunks {
		if _, err := outputFile.Write(chunk); err != nil {
			return fmt.Errorf("failed to write chunk %d: %w", i, err)
		}
	}

	fmt.Printf("File %s created successfully at %s\n", fileName, outputDir)
	return nil

}
