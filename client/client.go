package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	server "zerodupe/server/app"
)

// Client represents a client that can connect to a server and upload/download files
type Client struct {
	serverURL     string
	server        *server.Server
	serverStarted bool
	mu            sync.Mutex
}

// NewClient creates a new client
func NewClient(serverURL string) *Client {
	return &Client{
		serverURL: serverURL,
	}
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
func (client *Client) UploadFile(filePath string) error {
	content, contentHash, err := getFileHash(filePath)
	if err != nil {
		return err
	}

	reqBody := struct {
		FileHash   string `json:"filehash"`
		FileName   string `json:"file_name"`
		ChunkHash  string `json:"chunkhash"`
		ChunkOrder int    `json:"chunk_order"`
		Content    []byte `json:"content"`
	}{
		FileHash:   contentHash,
		FileName:   filePath,
		Content:    content,
		ChunkHash:  contentHash,
		ChunkOrder: 1,
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

	fmt.Printf("File uploaded successfully to %s\n", client.serverURL)
	return nil
}

// DownloadFile downloads a file from the server
func (client *Client) DownloadFile(fileName string) error {
	//TODO: how to handle getting multiple chunks
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

	err = os.WriteFile(fileName, fileContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("File downloaded successfully from %s\n", client.serverURL)
	return nil
}
