package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"zerodupe/internal/server/api"
	"zerodupe/internal/server/config"
)

func main() {
	port := flag.Int("port", 8080, "Server port")
	storageDir := flag.String("storage", "data/storage", "Storage directory")
	jwtSercet := flag.String("secret", "secret", "JWT Secret")
	accessTokenExpiry := flag.Int("Access Token Expiry", 15, "Access Token Expiry date in minutes")
	refreshTokenExpiry := flag.Int("Refresh Token Expiry", 24, "Refresh Token Expiry date in hours")
	flag.Parse()

	config := config.NewConfig(*port, *storageDir, *jwtSercet, *accessTokenExpiry, *refreshTokenExpiry)

	if err := os.MkdirAll(*storageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	fmt.Printf("Starting ZeroDupe server on port %d...\n", config.Port)
	fmt.Printf("Storage directory: %s\n", filepath.Clean(config.StorageDir))

	server.Run()
}
