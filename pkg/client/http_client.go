package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPClient implements the API interface using HTTP
type HTTPClient struct {
	serverURL  string
	httpClient *http.Client
	token      string
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(serverURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		serverURL: serverURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SetToken updates the client's authentication token
func (c *HTTPClient) SetToken(token string) {
	c.token = token
}

// Signup creates a new user account
func (c *HTTPClient) Signup(username, password, confirmPassword string) (error) {
	reqBody := SignUpRequest{
		Username:        username,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+"/auth/signup", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

// Login authenticates a user and returns access and refresh tokens
func (c *HTTPClient) Login(username, password string) (*AuthResponse, error) {
	reqBody := AuthRequest{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+"/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid credentials")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.SetToken(result.AccessToken)
	return &result, nil
}

// RefreshToken refreshes the access token using a refresh token
func (c *HTTPClient) RefreshToken(refreshToken string) (*AuthResponse, error) {
	reqBody := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+"/auth/refresh", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid refresh token")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.SetToken(result.AccessToken)
	return &result, nil
}
func (c *HTTPClient) addAuthHeader(req *http.Request) {
	if c.token != "" {
		token := c.token
		if !strings.HasPrefix(token, "Bearer ") {
			token = "Bearer " + token
		}
		req.Header.Add("Authorization", token)
	} else {
		fmt.Printf("DEBUG: No token available for request\n")
	}
}

// checkFileExists checks if a file exists on the server
func (c *HTTPClient) CheckFileExists(fileHash string) (bool, error) {
	req, err := http.NewRequest("GET", c.serverURL+"/check/"+fileHash, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return false, fmt.Errorf("unauthorized: invalid or missing token")
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("server error: %s", resp.Status)
	}

	var result FileExistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Exists, nil
}

// checkChunksExists checks if a list of chunks exists on the server and gets missing chunks from server
func (c *HTTPClient) GetMissingChunks(hashes []string) ([]string, error) {
	reqBody := MissingChunksRequest{
		Hashes: hashes,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+"/check", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid or missing token")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	var result MissingChunksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Missing, nil

}

// UploadChunk uploads a chunk to the server
func (c *HTTPClient) UploadChunk(request ChunkUploadRequest) (*ChunkUploadResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+"/upload", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	var result ChunkUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid or missing token")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	return &result, nil
}

// GetFileChunks gets the chunks hashes for a file from the server
func (c *HTTPClient) GetFileChunks(fileHash string) (*DownloadFileHashesResponse, error) {
	req, err := http.NewRequest("GET", c.serverURL+"/download/"+fileHash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid or missing token")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	var result DownloadFileHashesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DownloadChunk downloads a chunk's content from the server
func (c *HTTPClient) DownloadChunk(chunkHash string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.serverURL+"/chunk/"+chunkHash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid or missing token")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	chunkContent, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return chunkContent, nil

}
