package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"zerodupe/pkg/client"
)

func main() {
	// Authentication commands
	// zerodupe signup -server http://myhost:8080 -username user -password pass
	// zerodupe login -server http://myhost:8080 -username user -password pass
	// zerodupe refresh -server http://myhost:8080 -token <refresh_token>

	// File commands
	// zerodupe upload -server http://myhost:8080 -token <jwt_token> file.txt
	// zerodupe download -server http://myhost:8080 -token <jwt_token> -o ./downloads -n custom_name.txt HASH

	// Authentication command flags
	signupCmd := flag.NewFlagSet("signup", flag.ExitOnError)
	signupServerURL := signupCmd.String("server", "http://localhost:8080", "Server URL")
	signupUsername := signupCmd.String("username", "", "Username")
	signupPassword := signupCmd.String("password", "", "Password")

	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	loginServerURL := loginCmd.String("server", "http://localhost:8080", "Server URL")
	loginUsername := loginCmd.String("username", "", "Username")
	loginPassword := loginCmd.String("password", "", "Password")

	refreshCmd := flag.NewFlagSet("refresh", flag.ExitOnError)
	refreshServerURL := refreshCmd.String("server", "http://localhost:8080", "Server URL")
	refreshToken := refreshCmd.String("token", "", "Refresh token")

	// File command flags
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	uploadServerURL := uploadCmd.String("server", "http://localhost:8080", "Server URL")
	uploadToken := uploadCmd.String("token", "", "JWT authentication token")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadServerURL := downloadCmd.String("server", "http://localhost:8080", "Server URL")
	downloadToken := downloadCmd.String("token", "", "JWT authentication token")
	downloadOutput := downloadCmd.String("o", ".", "Output directory")
	downloadFileName := downloadCmd.String("n", "", "Output file name (default: file hash)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "signup":
		signupCmd.Parse(os.Args[2:])
		if *signupUsername == "" || *signupPassword == "" {
			fmt.Println("Error: Username and password are required")
			fmt.Println("Usage: zerodupe signup [flags]")
			signupCmd.PrintDefaults()
			os.Exit(1)
		}
		signup(*signupServerURL, *signupUsername, *signupPassword)

	case "login":
		loginCmd.Parse(os.Args[2:])
		if *loginUsername == "" || *loginPassword == "" {
			fmt.Println("Error: Username and password are required")
			fmt.Println("Usage: zerodupe login [flags]")
			loginCmd.PrintDefaults()
			os.Exit(1)
		}
		login(*loginServerURL, *loginUsername, *loginPassword)

	case "refresh":
		refreshCmd.Parse(os.Args[2:])
		if *refreshToken == "" {
			fmt.Println("Error: Refresh token is required")
			fmt.Println("Usage: zerodupe refresh [flags]")
			refreshCmd.PrintDefaults()
			os.Exit(1)
		}
		refresh(*refreshServerURL, *refreshToken)

	case "upload":
		uploadCmd.Parse(os.Args[2:])
		if uploadCmd.NArg() < 1 {
			fmt.Println("Error: Missing file path")
			fmt.Println("Usage: zerodupe upload [flags] <filepath>")
			uploadCmd.PrintDefaults()
			os.Exit(1)
		}
		filePath := uploadCmd.Arg(0)
		uploadFile(*uploadServerURL, *uploadToken, filePath)

	case "download":
		downloadCmd.Parse(os.Args[2:])
		if downloadCmd.NArg() < 1 {
			fmt.Println("Error: Missing file hash")
			fmt.Println("Usage: zerodupe download [flags] <filehash>")
			downloadCmd.PrintDefaults()
			os.Exit(1)
		}
		fileHash := downloadCmd.Arg(0)
		downloadFile(*downloadServerURL, *downloadToken, fileHash, *downloadOutput, *downloadFileName)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  zerodupe signup [flags]")
	fmt.Println("  zerodupe login [flags]")
	fmt.Println("  zerodupe refresh [flags]")
	fmt.Println("  zerodupe upload [flags] <filepath>")
	fmt.Println("  zerodupe download [flags] <filehash>")
	fmt.Println("\nFlags for signup:")
	fmt.Println("  -server string    Server URL (default \"http://localhost:8080\")")
	fmt.Println("  -username string  Username")
	fmt.Println("  -password string  Password")
	fmt.Println("\nFlags for login:")
	fmt.Println("  -server string    Server URL (default \"http://localhost:8080\")")
	fmt.Println("  -username string  Username")
	fmt.Println("  -password string  Password")
	fmt.Println("\nFlags for refresh:")
	fmt.Println("  -server string    Server URL (default \"http://localhost:8080\")")
	fmt.Println("  -token string     Refresh token")
	fmt.Println("\nFlags for upload:")
	fmt.Println("  -server string    Server URL (default \"http://localhost:8080\")")
	fmt.Println("  -token string     JWT authentication token")
	fmt.Println("\nFlags for download:")
	fmt.Println("  -server string    Server URL (default \"http://localhost:8080\")")
	fmt.Println("  -token string     JWT authentication token")
	fmt.Println("  -o string         Output directory (default \".\")")
	fmt.Println("  -n string         Output file name (default: file hash)")
}

// signup creates a new user account
func signup(serverURL, username, password string) {
	fmt.Printf("Creating account for user %s\n", username)

	client := client.NewClient(serverURL, "")
	resp, err := client.Signup(username, password)
	if err != nil {
		log.Fatalf("Failed to signup: %v", err)
	}

	fmt.Printf("Account created successfully\n")
	fmt.Printf("Access token: %s\n", resp.AccessToken)
	fmt.Printf("Refresh token: %s\n", resp.RefreshToken)
}

// login authenticates a user
func login(serverURL, username, password string) {
	fmt.Printf("Logging in user %s\n", username)

	client := client.NewClient(serverURL, "")
	resp, err := client.Login(username, password)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	fmt.Printf("Login successful\n")
	fmt.Printf("Access token: %s\n", resp.AccessToken)
	fmt.Printf("Refresh token: %s\n", resp.RefreshToken)
}

// refresh refreshes the access token
func refresh(serverURL, refreshToken string) {
	fmt.Println("Refreshing access token")

	client := client.NewClient(serverURL, "")
	resp, err := client.RefreshToken(refreshToken)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}

	fmt.Printf("Token refreshed successfully\n")
	fmt.Printf("New access token: %s\n", resp.AccessToken)
	fmt.Printf("New refresh token: %s\n", resp.RefreshToken)
}

// uploadFile uploads a file to the server
func uploadFile(serverURL, token, filePath string) {
	fmt.Printf("Uploading file %s to %s\n", filePath, serverURL)

	// Initialize client
	client := client.NewClient(serverURL, token)

	// Upload file
	err := client.UploadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
}

// downloadFile downloads a file from the server
func downloadFile(serverURL, token, fileHash, outputDir, fileName string) {
	fmt.Printf("Downloading file with hash %s from %s\n", fileHash, serverURL)

	client := client.NewClient(serverURL, token)

	err := client.DownloadFile(fileHash, outputDir, fileName)
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
}
