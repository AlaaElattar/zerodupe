package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"zerodupe/internal/server/auth"
	"zerodupe/internal/server/config"
	"zerodupe/internal/server/storage"
	"zerodupe/internal/server/storage/filesystem"
	"zerodupe/internal/server/storage/sqlite"

	_ "zerodupe/internal/server/docs" // This is the generated docs package

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server for all dependencies for server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     config.Config
	storage    storage.FileStorage
	handler    *Handler
}

// NewServer creates a new server with all configurations
func NewServer(config config.Config) (*Server, error) {
	fileStorage, err := filesystem.NewFilesystemStorage(config.StorageDir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create storage")
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	userStorage, err := sqlite.NewSqliteStorage(config.StorageDir + "/users.db")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user storage")
		return nil, fmt.Errorf("failed to create user storage: %w", err)
	}

	tokenHandler := auth.NewTokenHandler(
		config.JWTSecret,
		time.Duration(config.AccessTokenExpiryMin)*time.Minute,
		time.Duration(config.RefreshTokenExpiryHour)*time.Hour,
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

// registerHandlers registers all routes
func (server *Server) registerHandlers() {
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

// Run starts the server
func (server *Server) Run() error {
	addr := fmt.Sprintf(":%d", server.config.Port)
	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: server.router,
	}

	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Failed to start server")
	}

	return server.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (server *Server) Shutdown(ctx context.Context) error {
	if server.httpServer != nil {
		return server.httpServer.Shutdown(ctx)
	}
	return nil
}
