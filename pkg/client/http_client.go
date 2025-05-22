package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient implements the API interface using HTTP
type HTTPClient struct {
	serverURL  string
	httpClient *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(serverURL string, timeout time.Duration) *HTTPClient {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &HTTPClient{
		serverURL: serverURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// checkFileExists checks if a file exists on the server
func (c *HTTPClient) CheckFileExists(fileHash string) (bool, error) {
	resp, err := c.httpClient.Get(c.serverURL + "/check/" + fileHash)
	if err != nil {
		return false, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("server error: %s", resp.Status)
	}

	var result CheckFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Exists, nil
}

// checkChunksExists checks if a list of chunks exists on the server
func (c *HTTPClient) CheckChunksExists(hashes []string) ([]string, []string, error) {
	reqBody := CheckChunksRequest{
		Hashes: hashes,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := http.Post(c.serverURL+"/check", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("server error: %s", resp.Status)
	}

	var result CheckChunksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Exists, result.Missing, nil

}

// UploadChunk uploads a chunk to the server
func (c *HTTPClient) UploadChunk(request UploadRequest) (*UploadResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := c.httpClient.Post(c.serverURL+"/upload", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	var result UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	return &result, nil
}

func (c *HTTPClient) DownloadFileHashes(fileHash string) (*DownloadFileHashesResponse, error) {
	resp, err := c.httpClient.Get(c.serverURL + "/download/" + fileHash)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	var result DownloadFileHashesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *HTTPClient) DownloadChunkContent(chunkHash string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.serverURL + "/chunk/" + chunkHash)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	chunkContent, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return chunkContent, nil

}
