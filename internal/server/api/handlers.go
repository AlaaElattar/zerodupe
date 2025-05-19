package api

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"zerodupe/internal/server/model"
	"zerodupe/internal/server/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

// UploadFileHandler handles file upload requests
func (h *Handler) UploadFileHandler(c *gin.Context) {
	var request model.UploadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received chunk: Hash=%s, ChunkOrder=%d\n",
		request.FileHash, request.ChunkOrder)

	if err := h.storage.SaveChunkMetadata(request.FileHash, request.ChunkHash, request.ChunkOrder); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	hashMismatch := false
	if (len(request.Content)) > 0 {
		calculatedHash, err := h.storage.SaveChunkData(request.ChunkHash, request.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		hashMismatch = calculatedHash != request.ChunkHash
	}
	if (len(request.Content)) == 0 {
		exists, _, err := h.storage.CheckChunkExists([]string{request.ChunkHash})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if len(exists) == 0 {
			c.JSON(500, gin.H{"error": "Chunk does not exist"})
			return
		}
	}

	response := model.UploadResponse{
		Message:      "File uploaded successfully",
		FileHash:     request.FileHash,
		HashMismatch: hashMismatch,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CheckFileHashHandler(c *gin.Context) {
	fileHash := c.Param("filehash")

	if len(fileHash) < 4 {
		c.JSON(400, gin.H{"error": "Invalid file hash"})
		return
	}

	exists, err := h.storage.CheckFileExists(fileHash)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := model.CheckFileResponse{
		Exists: exists,
		Hash:   fileHash,
	}

	c.JSON(200, response)
}

func (h *Handler) CheckChunkHashesHandler(c *gin.Context) {
	var request struct {
		Hashes []string `json:"hashes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	exists, missing, err := h.storage.CheckChunkExists(request.Hashes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := model.CheckChunksResponse{
		Exists:  exists,
		Missing: missing,
	}

	c.JSON(200, response)
}

// DownloadFileHandler handles file download requests
func (h *Handler) DownloadFileHandler(c *gin.Context) {
	fileHash := c.Param("hash")

	fmt.Println("Downloading file with hash: " + fileHash)
	metadata, err := h.storage.GetFileMetadata(fileHash)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	sort.Slice(metadata.Chunks, func(i, j int) bool {
		return metadata.Chunks[i].ChunkOrder < metadata.Chunks[j].ChunkOrder
	})

	// Read each chunk and append to file content
	var fileContent []byte
	for _, chunk := range metadata.Chunks {
		chunkData, err := h.storage.GetChunkData(chunk.ChunkHash)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		fileContent = append(fileContent, chunkData...)
	}

	c.Data(http.StatusOK, "application/octet-stream", fileContent)

}
