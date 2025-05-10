package client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

// getFileHash returns the SHA-256 hash of a file's content
func getFileHash(filePath string) ([]byte, string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []byte{}, "", fmt.Errorf("file does not exist: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, "", err
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, "", err
	}

	fmt.Println(string(content))

	h := sha256.New()
	h.Write(content)
	contentHash := hex.EncodeToString(h.Sum(nil))
	fmt.Printf("File content SHA-256 hash: %s\n", contentHash)

	return content, contentHash, nil
}
