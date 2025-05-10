package server

import (
	"zerodupe/server/models"

	"github.com/gin-gonic/gin"
)

// Server for all dependencies for server
type Server struct {
	router *gin.Engine
	db     *models.DB
}

// NewServer creates a new server with all configurations
func NewServer() (*Server, error) {
	db, err := models.Connect("zerodupe.db")
	if err != nil {
		return nil, err
	}
	return &Server{
		router: gin.Default(),
		db:     db,
	}, nil
}

func (server *Server) registerHandlers() {
	server.router.POST("/upload", server.UploadFileHandler)
	server.router.GET("/download/:name", server.DownloadFileHandler)

}

func (server *Server) Run() {
	server.registerHandlers()
	server.router.Run(":8080")
}
