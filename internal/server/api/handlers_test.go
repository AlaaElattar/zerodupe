package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"zerodupe/internal/server/auth"
	"zerodupe/internal/server/model"
	"zerodupe/pkg/hasher"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupTestEnv sets up a test environment for the API handlers
func setupTestEnv() (*gin.Engine, *Handler, *MockFileStorage, *MockUserStorage, *MockTokenHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	fileStorage := new(MockFileStorage)
	userStorage := new(MockUserStorage)
	tokenHandler := new(MockTokenHandler)

	handler := NewHandler(fileStorage, userStorage, tokenHandler)

	return router, handler, fileStorage, userStorage, tokenHandler
}

func createTestFile(t *testing.T, name string, content []byte) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", name)
	require.NoError(t, err)

	defer tmpFile.Close()

	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	return tmpFile.Name()
}

func TestSignUpHandler(t *testing.T) {
	t.Run("Test SignUpHandler creates a new user", func(t *testing.T) {
		router, handler, _, userStorage, _ := setupTestEnv()
		router.POST("/auth/signup", handler.SignUpHandler)

		reqBody := model.AuthRequest{
			Username: "test",
			Password: "test",
		}

		userStorage.On("LoginUser", "test", "test").Return(nil, gorm.ErrRecordNotFound).Once()

		userStorage.On("CreateUser", mock.AnythingOfType("*model.User"), "test").Return(nil).Run(func(args mock.Arguments) {
			user := args.Get(0).(*model.User)
			assert.Equal(t, "test", user.Username)
		})

		userStorage.On("LoginUser", "test", "test").Return(&model.User{Username: "test"}, nil).Once()

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response auth.TokenPair
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)

		userStorage.AssertExpectations(t)

	})

	t.Run("Test SignUpHandler with Invalid request format", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.POST("/auth/signup", handler.SignUpHandler)

		req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer([]byte("invalid")))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")

	})

	t.Run("Test SignUpHandler with Username already exists", func(t *testing.T) {
		router, handler, _, userStorage, _ := setupTestEnv()
		router.POST("/auth/signup", handler.SignUpHandler)

		reqBody := model.AuthRequest{
			Username: "test",
			Password: "test",
		}

		userStorage.On("LoginUser", "test", "test").Return(&model.User{Username: "test"}, nil).Once()

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Username already exists")

		userStorage.AssertExpectations(t)

	})

	t.Run("Test SignUpHandler with Failed to create user", func(t *testing.T) {
		router, handler, _, userStorage, _ := setupTestEnv()
		router.POST("/auth/signup", handler.SignUpHandler)

		reqBody := model.AuthRequest{
			Username: "test",
			Password: "test",
		}

		userStorage.On("LoginUser", "test", "test").Return(nil, gorm.ErrRecordNotFound).Once()

		userStorage.On("CreateUser", mock.AnythingOfType("*model.User"), "test").Return(gorm.ErrDuplicatedKey).Once()

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create user")

		userStorage.AssertExpectations(t)

	})

	t.Run("Test SignUpHandler with Failed to login user", func(t *testing.T) {
		router, handler, _, userStorage, _ := setupTestEnv()
		router.POST("/auth/signup", handler.SignUpHandler)

		reqBody := model.AuthRequest{
			Username: "test",
			Password: "test",
		}

		userStorage.On("LoginUser", "test", "test").Return(nil, gorm.ErrRecordNotFound).Once()

		userStorage.On("CreateUser", mock.AnythingOfType("*model.User"), "test").Return(nil).Run(func(args mock.Arguments) {
			user := args.Get(0).(*model.User)
			assert.Equal(t, "test", user.Username)
		})

		userStorage.On("LoginUser", "test", "test").Return(nil, gorm.ErrRecordNotFound).Once()

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		userStorage.AssertExpectations(t)

	})
}

func TestLoginHandler(t *testing.T) {

	t.Run("Test LoginHandler with valid credentials", func(t *testing.T) {
		router, handler, _, userStorage, _ := setupTestEnv()
		router.POST("/auth/login", handler.LoginHandler)

		reqBody := model.AuthRequest{
			Username: "test",
			Password: "test",
		}

		userStorage.On("LoginUser", "test", "test").Return(&model.User{Username: "test"}, nil).Once()

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response auth.TokenPair
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)

		userStorage.AssertExpectations(t)

	})

	t.Run("Test LoginHandler with invalid request format", func(t *testing.T) {
		router, handler, _, userStorage, _ := setupTestEnv()
		router.POST("/auth/login", handler.LoginHandler)

		req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("invalid")))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")

		userStorage.AssertExpectations(t)

	})

	t.Run("Test LoginHandler with invalid credentials", func(t *testing.T) {
		router, handler, _, userStorage, _ := setupTestEnv()
		router.POST("/auth/login", handler.LoginHandler)

		reqBody := model.AuthRequest{
			Username: "test",
			Password: "test",
		}

		userStorage.On("LoginUser", "test", "test").Return(nil, gorm.ErrRecordNotFound).Once()

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")

		userStorage.AssertExpectations(t)

	})

}

func TestRefreshTokenHandler(t *testing.T) {

	t.Run("Test RefreshTokenHandler with valid refresh token", func(t *testing.T) {
		router, handler, _, _, tokenHandler := setupTestEnv()
		router.POST("/auth/refresh", handler.RefreshTokenHandler)

		claims := auth.TokenClaims{
			Username: "test",
			UserID:   1,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test-secret"))
		if err != nil {
			t.Errorf("Failed to sign token: %v", err)
		}

		reqBody := map[string]string{
			"refresh_token": tokenString,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			AccessToken string `json:"access_token"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
			return
		}

		assert.NotEmpty(t, response.AccessToken)

		_, err = tokenHandler.VerifyToken(response.AccessToken)
		if err != nil {
			t.Errorf("Failed to verify access token: %v", err)
		}

	})

	t.Run("Test RefreshTokenHandler with invalid request format", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.POST("/auth/refresh", handler.RefreshTokenHandler)

		req, err := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer([]byte("invalid")))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")

	})

	t.Run("Test RefreshTokenHandler with invalid refresh token", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.POST("/auth/refresh", handler.RefreshTokenHandler)

		reqBody := map[string]string{
			"refresh_token": "invalid",
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid refresh token")

	})

}

func TestUploadFileHandler(t *testing.T) {
	t.Run("Test UploadFileHandler with valid request (one chunk file)", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.POST("/upload", handler.UploadFileHandler)

		chunkHash := hasher.CalculateChunkHash([]byte("Hello World!"))
		content := []byte("Hello World!")

		uploadRequest := model.UploadRequest{
			FileHash:   chunkHash,
			ChunkHash:  chunkHash,
			ChunkOrder: 1,
			Content:    content,
		}

		fileStorage.On("SaveChunkData", chunkHash, content).Return(chunkHash, nil).Once()

		jsonBody, err := json.Marshal(uploadRequest)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UploadResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, "File uploaded successfully", response.Message)
		assert.Equal(t, chunkHash, response.FileHash)
		assert.False(t, response.HashMismatch)

		fileStorage.AssertExpectations(t)

	})

	t.Run("Test UploadFileHandler with valid request (multi chunk file)", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.POST("/upload", handler.UploadFileHandler)

		fileHash := "fileHash756456"
		chunkHash := hasher.CalculateChunkHash([]byte("Hello World!"))
		content := []byte("Hello World!")

		uploadRequest := model.UploadRequest{
			FileHash:   fileHash,
			ChunkHash:  chunkHash,
			ChunkOrder: 1,
			Content:    content,
		}

		fileStorage.On("SaveChunkMetadata", fileHash, chunkHash, 1).Return(nil).Once()
		fileStorage.On("SaveChunkData", chunkHash, content).Return(chunkHash, nil).Once()

		jsonBody, err := json.Marshal(uploadRequest)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UploadResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, "File uploaded successfully", response.Message)
		assert.Equal(t, fileHash, response.FileHash)
		assert.False(t, response.HashMismatch)

		fileStorage.AssertExpectations(t)

	})

	t.Run("Test UploadFileHandler with invalid request format", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.POST("/upload", handler.UploadFileHandler)

		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte("invalid")))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")

	})

	t.Run("Test UploadFileHandler with Internal Server Error", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.POST("/upload", handler.UploadFileHandler)

		chunkHash := hasher.CalculateChunkHash([]byte("Hello World!"))
		content := []byte("Hello World!")

		uploadRequest := model.UploadRequest{
			FileHash:   chunkHash,
			ChunkHash:  chunkHash,
			ChunkOrder: 1,
			Content:    content,
		}

		fileStorage.On("SaveChunkData", chunkHash, content).Return("", errors.New("test error")).Once()

		jsonBody, err := json.Marshal(uploadRequest)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to save chunk data")

		fileStorage.AssertExpectations(t)
	})

}

func TestCheckFileHashHandler(t *testing.T) {
	t.Run("Test CheckFileHashHandler with valid request", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.GET("/check/:filehash", handler.CheckFileHashHandler)

		fileHash := "fileHash756456"

		fileStorage.On("CheckFileExists", fileHash).Return(true, nil).Once()

		req, err := http.NewRequest("GET", "/check/"+fileHash, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.CheckFileResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.True(t, response.Exists)
		assert.Equal(t, fileHash, response.Hash)

		fileStorage.AssertExpectations(t)
	})

	t.Run("Test CheckFileHashHandler with invalid file hash", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.GET("/check/:filehash", handler.CheckFileHashHandler)

		fileHash := "fi"
		req, err := http.NewRequest("GET", "/check/"+fileHash, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid file hash")

	})

	t.Run("Test CheckFileHashHandler with internal server error", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.GET("/check/:filehash", handler.CheckFileHashHandler)

		fileHash := "fileHash756456"

		fileStorage.On("CheckFileExists", fileHash).Return(false, errors.New("test error")).Once()

		req, err := http.NewRequest("GET", "/check/"+fileHash, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "test error")

		fileStorage.AssertExpectations(t)
	})

}

func TestCheckChunkHashesHandler(t *testing.T) {
	t.Run("Test CheckChunkHashesHandler with valid request", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.POST("/check", handler.CheckChunkHashesHandler)

		hashes := []string{"hash1", "hash2", "hash3"}
		missing := []string{"hash2", "hash3"}

		fileStorage.On("CheckChunkExists", hashes).Return([]string{}, missing, nil).Once()

		reqBody := map[string][]string{
			"hashes": hashes,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/check", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.CheckChunksResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, missing, response.Missing)

		fileStorage.AssertExpectations(t)
	})

	t.Run("Test CheckChunkHashesHandler with invalid request format", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.POST("/check", handler.CheckChunkHashesHandler)

		req, err := http.NewRequest("POST", "/check", bytes.NewBuffer([]byte("invalid")))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

	t.Run("Test CheckChunkHashesHandler with internal server error", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.POST("/check", handler.CheckChunkHashesHandler)

		hashes := []string{"hash1", "hash2", "hash3"}

		fileStorage.On("CheckChunkExists", hashes).Return([]string{}, []string{}, errors.New("test error")).Once()

		reqBody := map[string][]string{
			"hashes": hashes,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/check", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "test error")

		fileStorage.AssertExpectations(t)
	})
}

func TestDownloadFileHandler(t *testing.T) {
	t.Run("Test DownloadFileHandler with valid request", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.GET("/download/:hash", handler.DownloadFileHandler)

		fileHash := "fileHash756456"
		fileStorage.On("GetFileMetadata", fileHash).Return(&model.FileMetadata{
			Chunks: []model.ChunkMetadata{
				{
					ChunkOrder: 1,
					ChunkHash:  "chunkhash123456789",
				},
			},
		}, nil).Once()

		req, err := http.NewRequest("GET", "/download/"+fileHash, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.DownloadFileResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, response.ChunkHashes, []string{"chunkhash123456789"})
		assert.Equal(t, response.ChunksCount, 1)
		assert.Equal(t, response.FileHash, fileHash)

	})

	t.Run("Test DownloadFileHandler with non-existing request", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.GET("/download/:hash", handler.DownloadFileHandler)

		req, _ := http.NewRequest("GET", "/download", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("Test DownloadFileHandler with internal server error", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.GET("/download/:hash", handler.DownloadFileHandler)

		fileHash := "fileHash756456"

		fileStorage.On("GetFileMetadata", fileHash).
			Return(&model.FileMetadata{}, errors.New("test metadata error")).Once()

		fileStorage.On("GetChunkData", fileHash).
			Return([]byte{}, errors.New("test chunk error")).Once()

		req, err := http.NewRequest("GET", "/download/"+fileHash, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "test chunk error")
	})

	t.Run("Test DownloadFileHandler with single chunk file", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.GET("/download/:hash", handler.DownloadFileHandler)

		fileHash := "fileHash756456"
		fileStorage.On("GetFileMetadata", fileHash).Return(&model.FileMetadata{}, errors.New("test error")).Once()
		fileStorage.On("GetChunkData", fileHash).
			Return([]byte("Hello World !"), nil).Once()

		req, err := http.NewRequest("GET", "/download/"+fileHash, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.DownloadFileResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, response.ChunkHashes, []string{fileHash})
		assert.Equal(t, response.ChunksCount, 1)
		assert.Equal(t, response.FileHash, fileHash)

	})

}

func TestGetChunkContent(t *testing.T) {

	t.Run("Test GetChunkContent with valid request", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.GET("/chunk/:hash", handler.GetChunkContent)

		chunkHash := "chunkHash1235487"
		fileStorage.On("GetChunkData", chunkHash).Return([]byte("Hello World!"), nil).Once()

		req, err := http.NewRequest("GET", "/chunk/"+chunkHash, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/octet-stream")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Test GetChunkContent with invalid request", func(t *testing.T) {
		router, handler, _, _, _ := setupTestEnv()
		router.GET("/chunk/:hash", handler.GetChunkContent)

		req, _ := http.NewRequest("GET", "/chunk/", nil)
		req.Header.Set("Content-Type", "application/octet-stream")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

	})
	t.Run("Test GetChunkContent with internal server error", func(t *testing.T) {
		router, handler, fileStorage, _, _ := setupTestEnv()
		router.GET("/chunk/:hash", handler.GetChunkContent)

		chunkHash := "chunkHash1235487"
		fileStorage.On("GetChunkData", chunkHash).Return([]byte(nil), errors.New("test error")).Once()

		req, err := http.NewRequest("GET", "/chunk/"+chunkHash, nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/octet-stream")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

}

func TestAuthMiddleware(t *testing.T) {

	t.Run("Test AuthMiddleware with valid request", func(t *testing.T) {
		router, _, _, _, tokenHandler := setupTestEnv()
		accessToken := "testaccesstoken123456789"
		refreshToken := "testrefreshtoken123456789"
		tokenHandler.On("CreateTokenPair", 0, "test").Return(&auth.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil).Once()

		router.Use(AuthMiddleware(tokenHandler))

		tokenHandler.On("VerifyToken", accessToken).Return(&auth.TokenClaims{
			Username: "test",
			UserID:   0,
		}, nil).Once()

		router.GET("/protected", func(c *gin.Context) {
			userID := c.GetInt("userID")
			username := c.GetString("username")
			c.JSON(200, gin.H{
				"msg":      "authorized",
				"userID":   userID,
				"username": username,
			})
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"msg": "authorized",
			"userID": 0,
			"username": "test"
		}`, w.Body.String())

	})

	t.Run("Test AuthMiddleware with no authorization header", func(t *testing.T) {
		router, _, _, _, tokenHandler := setupTestEnv()

		router.Use(AuthMiddleware(tokenHandler))

		router.GET("/protected", func(c *gin.Context) {
			userID := c.GetInt("userID")
			username := c.GetString("username")
			c.JSON(200, gin.H{
				"msg":      "authorized",
				"userID":   userID,
				"username": username,
			})
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

	})

	t.Run("Test AuthMiddleware with invalid token", func(t *testing.T) {
		router, _, _, _, tokenHandler := setupTestEnv()
		accessToken := "testaccesstoken123456789"
		refreshToken := "testrefreshtoken123456789"
		tokenHandler.On("CreateTokenPair", 0, "test").Return(&auth.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil).Once()

		router.Use(AuthMiddleware(tokenHandler))

		tokenHandler.On("VerifyToken", accessToken).Return(&auth.TokenClaims{}, errors.New("test error")).Once()

		router.GET("/protected", func(c *gin.Context) {
			userID := c.GetInt("userID")
			username := c.GetString("username")
			c.JSON(200, gin.H{
				"msg":      "authorized",
				"userID":   userID,
				"username": username,
			})
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

	})

}
