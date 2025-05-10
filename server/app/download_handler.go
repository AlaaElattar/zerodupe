package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DownloadFileHandler handles file download requests
func (server *Server) DownloadFileHandler(c *gin.Context) {
	name := c.Param("name")

	fmt.Println("Downloading file: " + name)

	fileInfo, err := server.db.GetFileByName(name)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	chunks, err := server.db.GetFileChunks(fileInfo.FileHash)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var fileContent []byte
	for _, chunk := range chunks {
		fileContent = append(fileContent, chunk.Chunk...)
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
}
