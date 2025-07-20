package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"zerodupe/internal/server/api"
	"zerodupe/internal/server/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serverConfig config.Config

var rootCmd = &cobra.Command{
	Use:   "Zerodupe Server",
	Short: "Deduplication file storage system that splits files into chunks and only stores unique chunks, saving storage space.",
	RunE: func(cmd *cobra.Command, args []string) error {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		// JWT Secret
		if serverConfig.JWTSecret == "" {
			serverConfig.JWTSecret = os.Getenv("JWT_SECRET")
			if serverConfig.JWTSecret == "" {
				log.Error().Msg("No JWT secret provided. Please set JWT_SECRET in your environment or config. Exiting.")
				os.Exit(1)
			}
		}

		// Storage Dir
		if serverConfig.StorageDir == "" {
			serverConfig.StorageDir = os.Getenv("STORAGE_DIR")
			if serverConfig.StorageDir == "" {
				serverConfig.StorageDir = "data/storage"
			}
		}

		// Port
		if serverConfig.Port == 0 {
			if portStr := os.Getenv("PORT"); portStr != "" {
				if port, err := strconv.Atoi(portStr); err == nil {
					serverConfig.Port = port
				}
			}
			if serverConfig.Port == 0 {
				serverConfig.Port = 8080
			}
		}

		// Access Token Expiry (minutes)
		if serverConfig.AccessTokenExpiryMin == 0 {
			if minStr := os.Getenv("ACCESS_TOKEN_EXPIRY_MIN"); minStr != "" {
				if min, err := strconv.Atoi(minStr); err == nil {
					serverConfig.AccessTokenExpiryMin = min
				}
			}
			if serverConfig.AccessTokenExpiryMin == 0 {
				serverConfig.AccessTokenExpiryMin = 30
			}
		}

		// Refresh Token Expiry (hours)
		if serverConfig.RefreshTokenExpiryHour == 0 {
			if hourStr := os.Getenv("REFRESH_TOKEN_EXPIRY_HOUR"); hourStr != "" {
				if hour, err := strconv.Atoi(hourStr); err == nil {
					serverConfig.RefreshTokenExpiryHour = hour
				}
			}
			if serverConfig.RefreshTokenExpiryHour == 0 {
				serverConfig.RefreshTokenExpiryHour = 24
			}
		}

		if err := os.MkdirAll(serverConfig.StorageDir, 0755); err != nil {
			log.Error().Err(err).Msg("Failed to create storage directory")
			return err
		}

		server, err := api.NewServer(serverConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating server")
		}

		// Graceful shutdown
		return gracefulShutdown(server)

	},
}

func gracefulShutdown(server *api.Server) error {
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
		return err
	}

	log.Info().Msg("Server gracefully stopped.")
	return nil

}

func init() {
	rootCmd.Flags().IntVarP(&serverConfig.Port, "port", "p", 8080, "Server port")
	rootCmd.Flags().StringVarP(&serverConfig.StorageDir, "storage", "s", "data/storage", "Storage directory")
	rootCmd.Flags().StringVarP(&serverConfig.JWTSecret, "secret", "", "", "JWT Secret")
	rootCmd.Flags().IntVar(&serverConfig.AccessTokenExpiryMin, "access-token-expiry-min", 30, "Access token expiry in minutes")
	rootCmd.Flags().IntVar(&serverConfig.RefreshTokenExpiryHour, "refresh-token-expiry-hour", 24, "Refresh token expiry in hours")
}

func Execute() error {
	return rootCmd.Execute()
}
