package client

import (
	"fmt"
	"strings"
	"time"
	"zerodupe/pkg/hasher"
)

// Client represents a client that uploads files to the server
type Client struct {
	api          API
	checker      *FileChecker
	uploader     *ChunkUploader
	downloader   *ChunkDownloader
	accessToken  string
	refreshToken string
}

// NewClient creates a new client
func NewClient(serverURL string) *Client {
	httpClient := NewHTTPClient(serverURL, 30*time.Minute)

	return &Client{
		api:        httpClient,
		checker:    NewFileChecker(httpClient),
		uploader:   NewUploader(httpClient),
		downloader: NewDownloader(httpClient),
	}
}

// SetTokens updates the client's authentication tokens
func (client *Client) SetTokens(accessToken, refreshToken string) {
	client.accessToken = accessToken
	client.refreshToken = refreshToken
	client.api.SetToken(accessToken)
}

// ExecuteWithAuth executes a function with authentication
func (client *Client) ExecuteWithAuth(fn func() error) error {
	err := fn()

	// If we get an auth error, try refreshing once and retry
	if err != nil && strings.Contains(err.Error(), "unauthorized") {
		if client.refreshToken == "" {
			return err
		}

		// Try to refresh the token
		resp, refreshErr := client.RefreshToken(client.refreshToken)
		if refreshErr != nil {
			return err
		}

		// Update tokens and retry
		client.SetTokens(resp.AccessToken, resp.RefreshToken)
		return fn()
	}

	return err
}

// UploadFile uploads a file to the server
// It handles file validation, chunking, and coordinating the upload process
func (client *Client) UploadFile(filePath string) error {
	// check if file exists
	if err := validateFile(filePath); err != nil {
		return err
	}

	// Read file content
	fileContent, err := getFileContent(filePath)
	if err != nil {
		return err
	}

	if len(fileContent) == 0 {
		return fmt.Errorf("cannot upload empty file")
	}

	// Split data into chunks and calculate file hash
	chunks, fileHash, err := hasher.SplitDataIntoChunks(fileContent)
	if err != nil {
		return err
	}

	// Check if file already exists on server
	exists, err := client.checker.CheckFileExists(fileHash)
	if err != nil {
		return err
	} else if exists {
		fmt.Printf("File already exists on server. Skipping upload.\n")
		fmt.Printf("File hash: %s\n", fileHash)
		return nil
	}

	fmt.Printf("File does not exist on server. Uploading...\n")
	fmt.Printf("File hash: %s\n", fileHash)
	fmt.Printf("Total chunks: %d\n", len(chunks))

	// Identify which chunks already exist on the server
	existingChunks, err := client.checker.IdentifyExistingChunks(chunks)
	if err != nil {
		return err
	}

	// make sever returns not existing chunks
	if err := client.uploader.UploadChunks(chunks, fileHash, filePath, existingChunks); err != nil {
		return err
	}

	fmt.Printf("File uploaded successfully\n")
	fmt.Printf("File hash: %s (use this hash to download the file)\n", fileHash)
	return nil

}

// DownloadFile downloads a file from the server
func (client *Client) DownloadFile(fileHash string, outputDir string, fileName string) error {
	exists, err := client.checker.CheckFileExists(fileHash)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("file does not exist on server")
	}

	fmt.Printf("File exists on server. Downloading...\n")

	var response *DownloadFileHashesResponse
	hashes, err := client.api.GetFileChunks(fileHash)
	if err != nil {
		return err
	}

	if hashes.ChunksCount == 0 {
		response = &DownloadFileHashesResponse{
			FileHash:    fileHash,
			ChunkHashes: []string{fileHash},
			ChunksCount: 1,
		}
	} else {
		response = hashes
	}

	content, err := client.downloader.DownloadChunks(response)
	if err != nil {
		return err
	}

	// combine chunks into file
	if err := hasher.CombineChunksIntoFile(content, outputDir, fileName); err != nil {
		return err
	}

	return nil
}

// Signup creates a new user account
func (client *Client) Signup(username, password, confirmPAssword string) (error) {
	return client.api.Signup(username, password, confirmPAssword)
}

// Login authenticates a user and returns access and refresh tokens
func (client *Client) Login(username, password string) (*AuthResponse, error) {
	return client.api.Login(username, password)
}

// RefreshToken refreshes the access token using a refresh token
func (client *Client) RefreshToken(refreshToken string) (*AuthResponse, error) {
	return client.api.RefreshToken(refreshToken)
}
