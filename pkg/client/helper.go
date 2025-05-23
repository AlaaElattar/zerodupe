package client

import (
	"fmt"
	"os"
)

// validateFile checks if a file exists or not
func validateFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	return nil
}

func getFileContent(filePath string) ([]byte, error) {
	// check if file exists
	if err := validateFile(filePath); err != nil {
		return nil, err
	}

	return os.ReadFile(filePath)
}
