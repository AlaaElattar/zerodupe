package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	server "zerodupe/server/app"
)

// Client represents a client that can connect to a server and upload/download files
type Client struct {
	serverURL     string
	server        *server.Server
	serverStarted bool
	mu            sync.Mutex
	uploadDir     string
	downloadDir   string
}

// NewClient creates a new client
func NewClient(serverURL string) *Client {
	return &Client{
		serverURL: serverURL,
	}
}

// SetUploadDir sets the directory for uploads
func (client *Client) SetUploadDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}
	client.uploadDir = dir
	return nil
}

// SetDownloadDir sets the directory for downloads
func (client *Client) SetDownloadDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}
	client.downloadDir = dir
	return nil
}

// StartServer starts the server if it's not already started
func (client *Client) StartServer() error {
	client.mu.Lock()
	defer client.mu.Unlock()

	if !client.serverStarted {
		server, err := server.NewServer()
		if err != nil {
			return err
		}
		client.server = server
		client.serverStarted = true
		go client.server.Run()
	}
	return nil
}

// UploadFile uploads a file to the server
func (client *Client) UploadFile(fileName string) error {
	filePath := filepath.Join(client.uploadDir, fileName)

	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	chunks, fileHash, err := client.getFileChunks(filePath)
	if err != nil {
		return err
	}

	fmt.Printf("File hash: %s\n", fileHash)
	fmt.Printf("Total chunks: %d\n", len(chunks))

	for i, chunk := range chunks {
		fmt.Printf("Uploading chunk %d/%d (Order: %d, Size: %d bytes, Hash: %s)\n",
			i+1, len(chunks), chunk.ChunkOrder, len(chunk.Data), chunk.ChunkHash)
		reqBody := struct {
			FileHash   string `json:"filehash"`
			FileName   string `json:"file_name"`
			ChunkHash  string `json:"chunkhash"`
			ChunkOrder int    `json:"chunk_order"`
			Content    []byte `json:"content"`
		}{
			FileHash:   fileHash,
			FileName:   fileName,
			Content:    chunk.Data,
			ChunkHash:  chunk.ChunkHash,
			ChunkOrder: chunk.ChunkOrder,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		resp, err := http.Post(client.serverURL+"/upload", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to connect to server: %w", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server error: %s", respBody)
		}

	}

	fmt.Printf("File uploaded successfully to %s\n", client.serverURL)
	return nil
}

// DownloadFile downloads a file from the server
func (client *Client) DownloadFile(fileName string) error {
	resp, err := http.Get(client.serverURL + "/download/" + fileName)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error: %s", resp.Status)
	}

	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	downloadPath := filepath.Join(client.downloadDir, fileName)

	err = os.WriteFile(downloadPath, fileContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("File downloaded successfully from %s\n", client.serverURL)
	return nil
}
