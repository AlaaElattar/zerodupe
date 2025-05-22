package client

import (
	"fmt"
	// "os"
	// "path/filepath"
	"time"
	"zerodupe/packages/hasher"
)

// Client represents a client that uploads files to the server
type Client struct {
	api        API
	checker    *FileChecker
	uploader   *Uploader
	downloader *Downloader
}

// NewClient creates a new client
func NewClient(serverURL string) *Client {
	httpClient := NewHTTPClient(serverURL, 30*time.Second)

	return &Client{
		api:        httpClient,
		checker:    NewFileChecker(httpClient),
		uploader:   NewUploader(httpClient),
		downloader: NewDownloader(httpClient),
	}
}

// UploadFile uploads a file to the server
func (client *Client) UploadFile(filePath string) error {
	// check if file exists
	if err := validateFile(filePath); err != nil {
		return err
	}

	// split file into chunks && get file hash
	chunks, fileHash, err := hasher.SplitFileIntoChunks(filePath)
	if err != nil {
		return err
	}

	// check file exists on server or not
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

	// check chunks exists on server or not
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

	hashes, err := client.api.DownloadFileHashes(fileHash)
	if err != nil {
		return err
	}

	content, err := client.downloader.DownloadChunks(hashes)
	if err != nil {
		return err
	}

	// combine chunks into file
	if err := hasher.CombineChunksIntoFile(content, outputDir, fileName); err != nil {
		return err
	}

	return nil
}
