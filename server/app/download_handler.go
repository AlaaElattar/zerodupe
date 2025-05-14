package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"net/http"

	"github.com/gin-gonic/gin"
)

// DownloadFileHandler handles file download requests
func (server *Server) DownloadFileHandler(c *gin.Context) {
	name := c.Param("name")

	fmt.Println("Downloading file: " + name)

	fileHash, err := server.db.GetFileHashByName(name)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	chunks, err := server.GetFileChunks(fileHash)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Sort the chunks by order
	chunkOrder := make([]int, 0, len(chunks))
	for order := range chunks {
		chunkOrder = append(chunkOrder, order)
	}
	sort.Ints(chunkOrder)

	// Read each chunk and append to file content
	var fileContent []byte
	for _, order := range chunkOrder {
		chunkHash := chunks[order]
		chunkPath := filepath.Join("server", "storage", "blocks", chunkHash[:4], chunkHash)
		chunkFile, err := os.OpenFile(chunkPath, os.O_RDONLY, 0644)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer chunkFile.Close()

		chunkData, err := io.ReadAll(chunkFile)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		fileContent = append(fileContent, chunkData...)
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Data(http.StatusOK, "application/octet-stream", fileContent)

}

func (server *Server) GetFileChunks(fileHash string) (map[int]string, error) {
	filePath := filepath.Join("server", "storage", "meta", fileHash[:4], fileHash)
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata FileMetadata
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &metadata)

	chunks := make(map[int]string)
	for _, chunk := range metadata.Chunks {
		chunks[chunk.ChunkOrder] = chunk.ChunkHash
		fmt.Println(chunk.ChunkOrder, chunk.ChunkHash)
	}

	return chunks, nil

}
