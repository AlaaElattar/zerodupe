package cmd

import (
	"context"
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
	"github.com/spf13/cobra"
)

var serverConfig config.Config

var rootCmd = &cobra.Command{
	Use:   "Zerodupe Server",
	Short: "Deduplication file storage system that splits files into chunks and only stores unique chunks, saving storage space.",
	RunE: func(cmd *cobra.Command, args []string) error {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		if serverConfig.JWTSecret == "" {
			log.Error().Msg("No JWT secret provided. Please set JWT_SECRET in your environment or config. Exiting.")
			os.Exit(1)
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
	rootCmd.Flags().IntVar(&serverConfig.AccessTokenExpiryMin, "access-token-expiry-min", 1, "Access token expiry in minutes")
	rootCmd.Flags().IntVar(&serverConfig.RefreshTokenExpiryHour, "refresh-token-expiry-hour", 24, "Refresh token expiry in hours")
}

func Execute() error {
	return rootCmd.Execute()
}
