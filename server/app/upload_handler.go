package server

import (
	"net/http"

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

// UploadFileHandler handles file upload requests
func (server *Server) UploadFileHandler(c *gin.Context) {
	var request UploadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save file to database
	if err := server.db.SaveFile(request.FileName, request.FileHash, request.ChunkHash, request.ChunkOrder, request.Content); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
