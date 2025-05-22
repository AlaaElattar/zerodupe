package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

// SplitFileIntoChunks reads a file and splits it into chunks
func SplitFileIntoChunks(filePath string) ([]FileChunk, string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, "", err
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	fileHasher := sha256.New()
	buffer := make([]byte, ChunkSizeBytes)
	chunkOrder := 1
	var chunks []FileChunk

	// Splitting the file into chunks
	for {
		bytes, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, "", err
		}
		if bytes == 0 {
			break
		}

		chunkData := make([]byte, bytes)
		copy(chunkData, buffer[:bytes])

		chunkHash := CalculateChunkHash(chunkData)
		fileHasher.Write([]byte(chunkHash))

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

	// Get the file hash
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
