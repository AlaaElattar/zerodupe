package client

import (
	"fmt"
	"os"
)

func validateFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	return nil
}
