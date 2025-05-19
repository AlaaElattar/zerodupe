package hasher

import (
	"os"
	"testing"
)

func createTestingFile(t *testing.T, name string, content []byte) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if err := os.WriteFile(tmpFile.Name(), content, 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestCalculateChunkHash(t *testing.T) {
	t.Run("Test CalculateChunkHash calculates correct hash for given data", func(t *testing.T) {
		data := []byte("hello world")
		expectedHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

		actualHash := CalculateChunkHash(data)
		if actualHash != expectedHash {
			t.Errorf("Expected hash %s, got %s", expectedHash, actualHash)
		}
	})

	t.Run("Test CalculateChunkHash calculates different hash for different data", func(t *testing.T) {
		data1 := []byte("hello world")
		data2 := []byte("hello world!")
		hash1 := CalculateChunkHash(data1)
		hash2 := CalculateChunkHash(data2)
		if hash1 == hash2 {
			t.Errorf("Expected different hashes for different data, got %s and %s", hash1, hash2)
		}
	})
}

func TestSplitFileIntoChunks(t *testing.T) {

	t.Run("Test SplitFileIntoChunks with empty file returns no chunks", func(t *testing.T) {
		filePath := createTestingFile(t, "empty.txt", []byte(""))
		chunks, _, err := SplitFileIntoChunks(filePath)
		expectedChunks := 0
		if err != nil {
			t.Fatalf("Failed to split file into chunks: %v", err)
		}
		if len(chunks) != expectedChunks {
			t.Errorf("Expected %d chunks, got %d", expectedChunks, len(chunks))
		}
	})

	t.Run("Test SplitFileIntoChunks with file smaller than chunk size returns one chunk", func(t *testing.T) {
		filePath := createTestingFile(t, "small.txt", []byte("hello world"))
		chunks, _, err := SplitFileIntoChunks(filePath)
		expectedChunks := 1
		if err != nil {
			t.Fatalf("Failed to split file into chunks: %v", err)
		}
		if len(chunks) != expectedChunks {
			t.Errorf("Expected %d chunks, got %d", expectedChunks, len(chunks))
		}
	})

	t.Run("Test SplitFileIntoChunks with file larger than chunk size returns multiple chunks", func(t *testing.T) {
		filaPath := createTestingFile(t, "big.txt", make([]byte, 3*ChunkSizeBytes))
		chunks, _, err := SplitFileIntoChunks(filaPath)
		expectedChunks := 3
		if err != nil {
			t.Fatalf("Failed to split file into chunks: %v", err)
		}
		if len(chunks) != expectedChunks {
			t.Errorf("Expected %d chunks, got %d", expectedChunks, len(chunks))
		}
	})
}

func TestVerifyChunkHash(t *testing.T) {
	t.Run("Test VerifyChunkHash returns true for correct hash", func(t *testing.T) {
		data := []byte("hello world")
		hash := CalculateChunkHash(data)
		isValid, _ := VerifyChunkHash(data, hash)
		if !isValid {
			t.Errorf("Expected valid hash, got invalid")
		}
	})

	t.Run("Test VerifyChunkHash returns false for incorrect hash", func(t *testing.T) {
		data := []byte("hello world")
		incorrectHash := "incorrectHash"
		isValid, _ := VerifyChunkHash(data, incorrectHash)
		if isValid {
			t.Errorf("Expected invalid hash, got valid")
		}
	})
}
