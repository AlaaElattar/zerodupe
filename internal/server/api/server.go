package api

import (
	"fmt"
	"time"
	"zerodupe/internal/server/auth"
	"zerodupe/internal/server/config"
	"zerodupe/internal/server/storage"
	"zerodupe/internal/server/storage/filesystem"
	"zerodupe/internal/server/storage/sqlite"

	"github.com/gin-gonic/gin"
)

// Server for all dependencies for server
type Server struct {
	router  *gin.Engine
	config  config.Config
	storage storage.FileStorage
	handler *Handler
}

// NewServer creates a new server with all configurations
func NewServer(config config.Config) (*Server, error) {
	fileStorage, err := filesystem.NewFilesystemStorage(config.StorageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	userStorage, err := sqlite.NewSqliteStorage(config.StorageDir + "/users.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create user storage: %w", err)
	}

	tokenHandler := auth.NewTokenHandler(
		config.JWTSecret,
		time.Duration(config.AccessTokenExpiry)*time.Minute,
		time.Duration(config.RefreshTokenExpiry)*time.Hour,
	)

	handler := NewHandler(fileStorage, userStorage, tokenHandler)
	router := gin.Default()

	server := &Server{
		router:  router,
		config:  config,
		handler: handler,
		storage: fileStorage,
	}

	// Register routes
	server.registerHandlers()

	return server, nil
}

func (server *Server) registerHandlers() {
	server.router.POST("/auth/signup", server.handler.SignUpHandler)
	server.router.POST("/auth/login", server.handler.LoginHandler)
	server.router.POST("/auth/refresh", server.handler.RefreshTokenHandler)

	authMiddleware := AuthMiddleware(server.handler.tokenHandler)
	authorized := server.router.Group("/")
	authorized.Use(authMiddleware)
	{
		authorized.POST("/upload", server.handler.UploadFileHandler)
		authorized.GET("/check/:filehash", server.handler.CheckFileHashHandler)
		authorized.POST("/check", server.handler.CheckChunkHashesHandler)
		authorized.GET("/download/:hash", server.handler.DownloadFileHandler)
		authorized.GET("/chunk/:hash", server.handler.GetChunkContent)
	}
}

func (server *Server) Run() {
	server.router.Run(fmt.Sprintf(":%d", server.config.Port))
}
