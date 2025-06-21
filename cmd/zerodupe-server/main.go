package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"zerodupe/internal/server/api"
	"zerodupe/internal/server/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var serverConfig config.Config

	flag.IntVar(&serverConfig.Port, "port", 8080, "Server port")
	flag.StringVar(&serverConfig.StorageDir, "storage", "data/storage", "Storage directory")
	flag.StringVar(&serverConfig.JWTSecret, "secret", "", "JWT Secret")
	flag.IntVar(&serverConfig.AccessTokenExpiry, "access-token-expiry", 15, "Access Token Expiry date in minutes")
	flag.IntVar(&serverConfig.RefreshTokenExpiry, "refresh-token-expiry", 24, "Refresh Token Expiry date in hours")
	flag.Parse()

	if serverConfig.JWTSecret == "" {
		// Generate a random secret if not provided
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		secret := make([]byte, 32)
		for i := range secret {
			secret[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
		serverConfig.JWTSecret = string(secret)
		log.Warn().Msg("No JWT secret provided. Generated a random secret. This is not secure for production use.")
	}

	if err := os.MkdirAll(serverConfig.StorageDir, 0755); err != nil {
		log.Error().Err(err).Msg("Failed to create storage directory")
	}

	server, err := api.NewServer(serverConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating server")
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Int("port", serverConfig.Port).Msg("Starting ZeroDupe server")
		log.Info().Str("path", filepath.Clean(serverConfig.StorageDir)).Msg("Storage directory")

		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Failed to start server")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Server gracefully stopped.")
}
