package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadRequest represents a file upload request
type UploadRequest struct {
	FileHash   string `json:"filehash" binding:"required"`
	FileName   string `json:"file_name" binding:"required"`
	ChunkHash  string `json:"chunkhash" binding:"required"`
	ChunkOrder int    `json:"chunk_order" binding:"required"`
	Content    []byte `json:"content" binding:"required"`
}

// ChunkMetadata represents metadata for a single chunk
type ChunkMetadata struct {
	ChunkOrder int    `json:"chunk_order"`
	ChunkHash  string `json:"chunk_hash"`
}

// FileMetadata represents metadata for a complete file
type FileMetadata struct {
	Chunks []ChunkMetadata `json:"chunks"`
}

// UploadFileHandler handles file upload requests
func (server *Server) UploadFileHandler(c *gin.Context) {
	var request UploadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received chunk: File=%s, Hash=%s, ChunkOrder=%d\n",
		request.FileName, request.FileHash, request.ChunkOrder)

	if err := server.checkStorageDir(request.FileHash, request.ChunkHash); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := server.saveChunkMetadata(request.FileHash, request.ChunkHash, request.ChunkOrder); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := server.saveChunkData(request.ChunkHash, request.Content); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := server.db.SaveFile(request.FileName, request.FileHash); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Directories created successfully"})
}

// checkStorageDir checks if a directory exists and creates it if it doesn't
func (server *Server) checkStorageDir(fileHash string, chunkHash string) error {

	directories := []string{
		filepath.Join("server", "storage"),
		filepath.Join("server", "storage", "blocks"),
		filepath.Join("server", "storage", "meta"),
		filepath.Join("server", "storage", "meta", fileHash[:4]),
		filepath.Join("server", "storage", "blocks", chunkHash[:4]),
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
	metaPath := filepath.Join("server", "storage", "meta", fileHash[:4], fileHash)

	if err := os.MkdirAll(filepath.Dir(metaPath), os.ModePerm); err != nil {
		return err
	}

	var metadata FileMetadata
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

	metadata.Chunks = append(metadata.Chunks, ChunkMetadata{
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
func (server *Server) saveChunkData(chunkHash string, content []byte) error {
	blockPath := filepath.Join("server", "storage", "blocks", chunkHash[:4], chunkHash)

	if _, err := os.Stat(blockPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(blockPath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(blockPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}

	return nil
}
