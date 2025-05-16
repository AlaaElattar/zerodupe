package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"zerodupe/internal/server/model"
	"zerodupe/packages/hasher"
)

func (server *Server) CheckFileExists(fileHash string) (bool, error) {
	metaPath := filepath.Join(server.config.StorageDir, "meta", fileHash[:4], fileHash)
	if _, err := os.Stat(metaPath); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}

	return false, nil
}

func (server *Server) CheckChunkExists(hashes []string) ([]string, []string, error) {
	var existingChunks []string
	var missingChunks []string

	for _, hash := range hashes {
		blockPath := filepath.Join(server.config.StorageDir, "blocks", hash[:4], hash)

		_, err := os.Stat(blockPath)
		if err != nil {
			if os.IsNotExist(err) {
				missingChunks = append(missingChunks, hash)
			} else {
				return nil, nil, err
			}
		} else {
			existingChunks = append(existingChunks, hash)
		}
	}

	return existingChunks, missingChunks, nil
}

// checkStorageDir checks if a directory exists and creates it if it doesn't
func (server *Server) checkStorageDir(fileHash string, chunkHash string) error {

	directories := []string{
		server.config.StorageDir,
		filepath.Join(server.config.StorageDir, "blocks"),
		filepath.Join(server.config.StorageDir, "meta"),
		filepath.Join(server.config.StorageDir, "meta", fileHash[:4]),
		filepath.Join(server.config.StorageDir, "blocks", chunkHash[:4]),
	}

	for _, dirPath := range directories {
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

// saveChunkMetadata saves file metadata to a file
func (server *Server) saveChunkMetadata(fileHash string, chunkHash string, chunkOrder int) error {
	metaPath := filepath.Join(server.config.StorageDir, "meta", fileHash[:4], fileHash)

	if err := os.MkdirAll(filepath.Dir(metaPath), os.ModePerm); err != nil {
		return err
	}

	var metadata model.FileMetadata
	if _, err := os.Stat(metaPath); err == nil {
		content, err := os.ReadFile(metaPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(content, &metadata); err != nil {
			return err
		}
	}

	for _, chunk := range metadata.Chunks {
		if chunk.ChunkOrder == chunkOrder && chunk.ChunkHash == chunkHash {
			return nil
		}
	}

	metadata.Chunks = append(metadata.Chunks, model.ChunkMetadataRequest{
		ChunkOrder: chunkOrder,
		ChunkHash:  chunkHash,
	})

	newContent, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	if err := os.WriteFile(metaPath, newContent, 0644); err != nil {
		return err
	}

	return nil
}

// saveChunkData saves chunk data to a file
func (server *Server) saveChunkData(chunkHash string, content []byte) (string, error) {
	blockPath := filepath.Join(server.config.StorageDir, "blocks", chunkHash[:4], chunkHash)

	if _, err := os.Stat(blockPath); err == nil {
		return "", nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(blockPath), os.ModePerm); err != nil {
		return "", err
	}

	isValid, calculatedHash := hasher.VerifyChunkHash(content, chunkHash)
	if !isValid {
		fmt.Printf("WARNING: Hash mismatch. Expected: %s, Got: %s\n", chunkHash, calculatedHash)
	}

	file, err := os.OpenFile(blockPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return calculatedHash, err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return calculatedHash, err
	}

	return calculatedHash, nil
}
