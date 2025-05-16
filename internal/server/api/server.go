package api

import (
	"fmt"
	"zerodupe/internal/server/config"
	"zerodupe/internal/server/storage"
	"zerodupe/internal/server/storage/filesystem"

	"github.com/gin-gonic/gin"
)

// Server for all dependencies for server
type Server struct {
	router  *gin.Engine
	config  config.Config
	storage storage.Storage
	handler *Handler
}

// NewServer creates a new server with all configurations
func NewServer(config config.Config) (*Server, error) {
	storage, err := filesystem.NewFilesystemStorage(config.StorageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}
	handler := NewHandler(storage)
	router := gin.Default()

	server := &Server{
		router:  router,
		config:  config,
		handler: handler,
		storage: storage,
	}

	// Register routes
	server.registerHandlers()

	return server, nil
}

func (server *Server) registerHandlers() {
	server.router.POST("/upload", server.handler.UploadFileHandler)
	server.router.GET("/check/:filehash", server.handler.CheckFileHashHandler)
	server.router.POST("/check", server.handler.CheckChunkHashesHandler)
	server.router.GET("/download/:hash", server.handler.DownloadFileHandler)

}

func (server *Server) Run() {
	server.router.Run(fmt.Sprintf(":%d", server.config.Port))
}
