package api

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"
// 	"zerodupe/internal/server/auth"
// 	"zerodupe/internal/server/model"
// 	"zerodupe/pkg/hasher"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// 	"gorm.io/gorm"
// )

// // setupTestEnv sets up a test environment for the API handlers
// func setupTestEnv() (*gin.Engine, *Handler, *MockFileStorage, *MockUserStorage, *MockTokenHandler) {
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()

// 	fileStorage := new(MockFileStorage)
// 	userStorage := new(MockUserStorage)
// 	tokenHandler := new(MockTokenHandler)

// 	handler := NewHandler(fileStorage, userStorage, tokenHandler)

// 	return router, handler, fileStorage, userStorage, tokenHandler
// }

// // create a new HTTP request with a JSON body and set Content-Type
// func newRequest(t *testing.T, method, url string, body interface{}) *http.Request {
// 	var buf bytes.Buffer
// 	if body != nil {
// 		err := json.NewEncoder(&buf).Encode(body)
// 		require.NoError(t, err)
// 	}
// 	req, err := http.NewRequest(method, url, &buf)
// 	require.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")
// 	return req
// }

// // perform a request and return the response
// func applyRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	return w
// }

// func Test_SignUpHandler(t *testing.T) {
// 	t.Run("Test_SignUpHandler_Creates_A_New_User", func(t *testing.T) {
// 		router, handler, _, userStorage, _ := setupTestEnv()
// 		router.POST("/auth/signup", handler.SignUpHandler)

// 		username := "username"

// 		reqBody := SignUpRequest{
// 			Username:        username,
// 			Password:        "test",
// 			ConfirmPassword: "test",
// 		}

// 		userStorage.On("GetUserByUsername", username).Return(nil, gorm.ErrRecordNotFound).Once()

// 		userStorage.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
// 			user := args.Get(0).(*model.User)
// 			assert.Equal(t, username, user.Username)
// 		})

// 		req := newRequest(t, "POST", "/auth/signup", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.JSONEq(t, `"user registered successfully"`, w.Body.String())
// 		userStorage.AssertExpectations(t)
// 	})

// 	t.Run("Test_SignUpHandler_With_Invalid_Request_Format", func(t *testing.T) {
// 		router, handler, _, _, _ := setupTestEnv()
// 		router.POST("/auth/signup", handler.SignUpHandler)

// 		req := newRequest(t, "POST", "/auth/signup", nil)
// 		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid"))) // override with invalid body
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		assert.Contains(t, w.Body.String(), "Invalid request format")
// 	})

// 	t.Run("Test_SignUpHandler_With_Username_Already_Exists", func(t *testing.T) {
// 		router, handler, _, userStorage, _ := setupTestEnv()
// 		router.POST("/auth/signup", handler.SignUpHandler)
// 		username := "username"

// 		reqBody := SignUpRequest{
// 			Username:        username,
// 			Password:        "test",
// 			ConfirmPassword: "test",
// 		}

// 		userStorage.On("GetUserByUsername", username).
// 			Return(&model.User{Username: username}, nil).
// 			Once()

// 		req := newRequest(t, "POST", "/auth/signup", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusConflict, w.Code)
// 		assert.Contains(t, w.Body.String(), "user already exists")

// 		userStorage.AssertExpectations(t)

// 	})

// 	t.Run("Test_SignUpHandler_With_Failed_To_Create_User", func(t *testing.T) {
// 		router, handler, _, userStorage, _ := setupTestEnv()
// 		router.POST("/auth/signup", handler.SignUpHandler)

// 		reqBody := SignUpRequest{
// 			Username:        "test",
// 			Password:        "test",
// 			ConfirmPassword: "test",
// 		}

// 		userStorage.On("GetUserByUsername", "test").Return(nil, gorm.ErrRecordNotFound).Once()

// 		userStorage.On("CreateUser", mock.AnythingOfType("*model.User")).Return(gorm.ErrDuplicatedKey).Once()

// 		req := newRequest(t, "POST", "/auth/signup", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		assert.Contains(t, w.Body.String(), "internal server error")

// 		userStorage.AssertExpectations(t)

// 	})
// }

// func Test_LoginHandler(t *testing.T) {

// 	t.Run("Test_LoginHandler_With_Valid_Credentials", func(t *testing.T) {
// 		router, handler, _, userStorage, tokenHandler := setupTestEnv()
// 		router.POST("/auth/login", handler.LoginHandler)

// 		reqBody := LoginRequest{
// 			Username: "test",
// 			Password: "test",
// 		}

// 		user := &model.User{
// 			ID:       1,
// 			Username: "test",
// 			Password: []byte("hashed"),
// 		}

// 		userStorage.On("GetUserByUsername", "test").Return(user, nil).Once()

// 		tokenHandler.On("CreateTokenPair", uint(1), user.Username).Return(&auth.TokenPair{
// 			AccessToken:  "access-token",
// 			RefreshToken: "refresh-token",
// 		}, nil).Once()

// 		req := newRequest(t, "POST", "/auth/login", reqBody)
// 		w := applyRequest(router, req)
// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response auth.TokenPair
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)

// 		assert.Equal(t, "access-token", response.AccessToken)
// 		assert.Equal(t, "refresh-token", response.RefreshToken)

// 		userStorage.AssertExpectations(t)
// 		tokenHandler.AssertExpectations(t)

// 	})

// 	t.Run("Test_LoginHandler_With_Invalid_Request_Format", func(t *testing.T) {
// 		router, handler, _, userStorage, _ := setupTestEnv()
// 		router.POST("/auth/login", handler.LoginHandler)

// 		req := newRequest(t, "POST", "/auth/login", nil)
// 		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid"))) // override with invalid body
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		assert.Contains(t, w.Body.String(), "Invalid request format")

// 		userStorage.AssertExpectations(t)

// 	})

// 	t.Run("Test_LoginHandler_With_Non_Existing_User", func(t *testing.T) {
// 		router, handler, _, userStorage, _ := setupTestEnv()
// 		router.POST("/auth/login", handler.LoginHandler)

// 		reqBody := LoginRequest{
// 			Username: "test",
// 			Password: "test",
// 		}

// 		userStorage.On("GetUserByUsername", "test").Return(nil, gorm.ErrRecordNotFound).Once()

// 		req := newRequest(t, "POST", "/auth/login", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusNotFound, w.Code)
// 		assert.Contains(t, w.Body.String(), "user does not exist")

// 		userStorage.AssertExpectations(t)

// 	})

// }

// func Test_RefreshTokenHandler(t *testing.T) {

// 	t.Run("Test_RefreshTokenHandler_With_Valid_Refresh_Token", func(t *testing.T) {
// 		router, handler, _, _, tokenHandler := setupTestEnv()
// 		router.POST("/auth/refresh", handler.RefreshTokenHandler)

// 		claims := auth.TokenClaims{
// 			Username: "test",
// 			UserID:   1,
// 			StandardClaims: jwt.StandardClaims{
// 				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
// 				IssuedAt:  time.Now().Unix(),
// 			},
// 		}
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 		tokenString, err := token.SignedString([]byte("test-secret"))
// 		assert.NoError(t, err)

// 		tokenHandler.On("RefreshAccessToken", tokenString).Return("new-access-token", nil).Once()

// 		reqBody := map[string]string{
// 			"refresh_token": tokenString,
// 		}

// 		req := newRequest(t, "POST", "/auth/refresh", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response struct {
// 			AccessToken string `json:"access_token"`
// 		}
// 		err = json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)

// 		assert.NotEmpty(t, response.AccessToken)

// 	})

// 	t.Run("Test_RefreshTokenHandler_With_Invalid_Request_Format", func(t *testing.T) {
// 		router, handler, _, _, _ := setupTestEnv()
// 		router.POST("/auth/refresh", handler.RefreshTokenHandler)

// 		req := newRequest(t, "POST", "/auth/refresh", nil)
// 		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid"))) // override with invalid body
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		assert.Contains(t, w.Body.String(), "Invalid request format")

// 	})

// 	t.Run("Test_RefreshTokenHandler_With_Invalid_Refresh_Token", func(t *testing.T) {
// 		router, handler, _, _, tokenHandler := setupTestEnv()
// 		router.POST("/auth/refresh", handler.RefreshTokenHandler)

// 		reqBody := map[string]string{
// 			"refresh_token": "invalid",
// 		}
// 		tokenHandler.On("RefreshAccessToken", "invalid").Return("", errors.New("invalid token")).Once()
// 		req := newRequest(t, "POST", "/auth/refresh", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusUnauthorized, w.Code)
// 		assert.Contains(t, w.Body.String(), "Invalid refresh token")

// 	})

// }

// func Test_UploadFileHandler(t *testing.T) {
// 	t.Run("Test_UploadFileHandler_With_Valid_Request_One_Chunk_File", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.POST("/upload", handler.UploadFileHandler)

// 		chunkHash := hasher.CalculateChunkHash([]byte("Hello World!"))
// 		content := []byte("Hello World!")

// 		uploadRequest := UploadRequest{
// 			FileHash:   chunkHash,
// 			ChunkHash:  chunkHash,
// 			ChunkOrder: 1,
// 			Content:    content,
// 		}

// 		fileStorage.On("SaveChunkData", chunkHash, content).Return(chunkHash, nil).Once()

// 		req := newRequest(t, "POST", "/upload", uploadRequest)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response UploadResponse
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		require.NoError(t, err)

// 		assert.Equal(t, "File uploaded successfully", response.Message)
// 		assert.Equal(t, chunkHash, response.FileHash)
// 		assert.False(t, response.HashMismatch)

// 		fileStorage.AssertExpectations(t)

// 	})

// 	t.Run("Test_UploadFileHandler_With_Valid_Request_Multi_Chunk_File", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.POST("/upload", handler.UploadFileHandler)

// 		fileHash := "fileHash756456"
// 		chunkHash := hasher.CalculateChunkHash([]byte("Hello World!"))
// 		content := []byte("Hello World!")

// 		uploadRequest := UploadRequest{
// 			FileHash:   fileHash,
// 			ChunkHash:  chunkHash,
// 			ChunkOrder: 1,
// 			Content:    content,
// 		}

// 		fileStorage.On("SaveChunkMetadata", fileHash, chunkHash, 1).Return(nil).Once()
// 		fileStorage.On("SaveChunkData", chunkHash, content).Return(chunkHash, nil).Once()

// 		req := newRequest(t, "POST", "/upload", uploadRequest)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response UploadResponse
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		require.NoError(t, err)

// 		assert.Equal(t, "File uploaded successfully", response.Message)
// 		assert.Equal(t, fileHash, response.FileHash)
// 		assert.False(t, response.HashMismatch)

// 		fileStorage.AssertExpectations(t)

// 	})

// 	t.Run("Test_UploadFileHandler_With_Invalid_Request_Format", func(t *testing.T) {
// 		router, handler, _, _, _ := setupTestEnv()
// 		router.POST("/upload", handler.UploadFileHandler)

// 		req := newRequest(t, "POST", "/upload", nil)
// 		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid"))) // override with invalid body
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		assert.Contains(t, w.Body.String(), "Invalid request format")

// 	})

// 	t.Run("Test_UploadFileHandler_With_Internal_Server_Error", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.POST("/upload", handler.UploadFileHandler)

// 		chunkHash := hasher.CalculateChunkHash([]byte("Hello World!"))
// 		content := []byte("Hello World!")

// 		uploadRequest := UploadRequest{
// 			FileHash:   chunkHash,
// 			ChunkHash:  chunkHash,
// 			ChunkOrder: 1,
// 			Content:    content,
// 		}

// 		fileStorage.On("SaveChunkData", chunkHash, content).Return("", errors.New("test error")).Once()

// 		req := newRequest(t, "POST", "/upload", uploadRequest)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		assert.Contains(t, w.Body.String(), "Failed to save chunk data")

// 		fileStorage.AssertExpectations(t)
// 	})

// }

// func Test_CheckFileHashHandler(t *testing.T) {
// 	t.Run("Test_CheckFileHashHandler_With_Valid_Request", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.GET("/check/:filehash", handler.CheckFileHashHandler)

// 		fileHash := "fileHash756456"

// 		fileStorage.On("CheckFileExists", fileHash).Return(true, nil).Once()

// 		req := newRequest(t, "GET", "/check/"+fileHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response CheckFileResponse
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		require.NoError(t, err)

// 		assert.True(t, response.Exists)
// 		assert.Equal(t, fileHash, response.Hash)

// 		fileStorage.AssertExpectations(t)
// 	})

// 	t.Run("Test_CheckFileHashHandler_With_Invalid_File_Hash", func(t *testing.T) {
// 		router, handler, _, _, _ := setupTestEnv()
// 		router.GET("/check/:filehash", handler.CheckFileHashHandler)

// 		fileHash := "fi"
// 		req := newRequest(t, "GET", "/check/"+fileHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		assert.Contains(t, w.Body.String(), "Invalid file hash")

// 	})

// 	t.Run("Test_CheckFileHashHandler_With_Internal_Server_Error", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.GET("/check/:filehash", handler.CheckFileHashHandler)

// 		fileHash := "fileHash756456"

// 		fileStorage.On("CheckFileExists", fileHash).Return(false, errors.New("test error")).Once()

// 		req := newRequest(t, "GET", "/check/"+fileHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		assert.Contains(t, w.Body.String(), "test error")

// 		fileStorage.AssertExpectations(t)
// 	})

// }

// func Test_CheckChunkHashesHandler(t *testing.T) {
// 	t.Run("Test_CheckChunkHashesHandler_With_Valid_Request", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.POST("/check", handler.CheckChunkHashesHandler)

// 		hashes := []string{"hash1", "hash2", "hash3"}
// 		missing := []string{"hash2", "hash3"}

// 		fileStorage.On("CheckChunkExists", hashes).Return([]string{}, missing, nil).Once()

// 		reqBody := map[string][]string{
// 			"hashes": hashes,
// 		}

// 		req := newRequest(t, "POST", "/check", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response CheckChunksResponse
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		require.NoError(t, err)

// 		assert.Equal(t, missing, response.Missing)

// 		fileStorage.AssertExpectations(t)
// 	})

// 	t.Run("Test_CheckChunkHashesHandler_With_Invalid_Request_Format", func(t *testing.T) {
// 		router, handler, _, _, _ := setupTestEnv()
// 		router.POST("/check", handler.CheckChunkHashesHandler)

// 		req := newRequest(t, "POST", "/check", nil)
// 		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid"))) // override with invalid body
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 	})

// 	t.Run("Test_CheckChunkHashesHandler_With_Internal_Server_Error", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.POST("/check", handler.CheckChunkHashesHandler)

// 		hashes := []string{"hash1", "hash2", "hash3"}

// 		fileStorage.On("CheckChunkExists", hashes).Return([]string{}, []string{}, errors.New("test error")).Once()

// 		reqBody := map[string][]string{
// 			"hashes": hashes,
// 		}

// 		req := newRequest(t, "POST", "/check", reqBody)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		assert.Contains(t, w.Body.String(), "test error")

// 		fileStorage.AssertExpectations(t)
// 	})
// }

// func Test_DownloadFileHandler(t *testing.T) {
// 	t.Run("Test_DownloadFileHandler_With_Valid_Request", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.GET("/download/:hash", handler.DownloadFileHandler)

// 		fileHash := "fileHash756456"
// 		fileStorage.On("GetFileMetadata", fileHash).Return(&model.FileMetadata{
// 			Chunks: []model.ChunkMetadata{
// 				{
// 					ChunkOrder: 1,
// 					ChunkHash:  "chunkhash123456789",
// 				},
// 			},
// 		}, nil).Once()

// 		req := newRequest(t, "GET", "/download/"+fileHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response DownloadFileResponse
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		require.NoError(t, err)

// 		assert.Equal(t, response.ChunkHashes, []string{"chunkhash123456789"})
// 		assert.Equal(t, response.ChunksCount, 1)
// 		assert.Equal(t, response.FileHash, fileHash)

// 	})

// 	t.Run("Test_DownloadFileHandler_With_Non_Existing_Request", func(t *testing.T) {
// 		router, handler, _, _, _ := setupTestEnv()
// 		router.GET("/download/:hash", handler.DownloadFileHandler)

// 		req := newRequest(t, "GET", "/download", nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusNotFound, w.Code)

// 	})

// 	t.Run("Test_DownloadFileHandler_With_Internal_Server_Error", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.GET("/download/:hash", handler.DownloadFileHandler)

// 		fileHash := "fileHash756456"

// 		fileStorage.On("GetFileMetadata", fileHash).
// 			Return(&model.FileMetadata{}, errors.New("test metadata error")).Once()

// 		fileStorage.On("GetChunkData", fileHash).
// 			Return([]byte{}, errors.New("test chunk error")).Once()

// 		req := newRequest(t, "GET", "/download/"+fileHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		assert.Contains(t, w.Body.String(), "test chunk error")
// 	})

// 	t.Run("Test_DownloadFileHandler_With_Single_Chunk_File", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.GET("/download/:hash", handler.DownloadFileHandler)

// 		fileHash := "fileHash756456"
// 		fileStorage.On("GetFileMetadata", fileHash).Return(&model.FileMetadata{}, errors.New("test error")).Once()
// 		fileStorage.On("GetChunkData", fileHash).
// 			Return([]byte("Hello World !"), nil).Once()

// 		req := newRequest(t, "GET", "/download/"+fileHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		var response DownloadFileResponse
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		require.NoError(t, err)

// 		assert.Equal(t, response.ChunkHashes, []string{fileHash})
// 		assert.Equal(t, response.ChunksCount, 1)
// 		assert.Equal(t, response.FileHash, fileHash)

// 	})

// }

// func Test_GetChunkContent(t *testing.T) {

// 	t.Run("Test_GetChunkContent_With_Valid_Request", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.GET("/chunk/:hash", handler.GetChunkContent)

// 		chunkHash := "chunkHash1235487"
// 		fileStorage.On("GetChunkData", chunkHash).Return([]byte("Hello World!"), nil).Once()

// 		req := newRequest(t, "GET", "/chunk/"+chunkHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 	})

// 	t.Run("Test_GetChunkContent_With_Invalid_Request", func(t *testing.T) {
// 		router, handler, _, _, _ := setupTestEnv()
// 		router.GET("/chunk/:hash", handler.GetChunkContent)

// 		req := newRequest(t, "GET", "/chunk/", nil)
// 		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid"))) // override with invalid body
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusNotFound, w.Code)

// 	})
// 	t.Run("Test_GetChunkContent_With_Internal_Server_Error", func(t *testing.T) {
// 		router, handler, fileStorage, _, _ := setupTestEnv()
// 		router.GET("/chunk/:hash", handler.GetChunkContent)

// 		chunkHash := "chunkHash1235487"
// 		fileStorage.On("GetChunkData", chunkHash).Return([]byte(nil), errors.New("test error")).Once()

// 		req := newRequest(t, "GET", "/chunk/"+chunkHash, nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)

// 	})

// }

// func Test_AuthMiddleware(t *testing.T) {

// 	t.Run("Test_AuthMiddleware_With_Valid_Request", func(t *testing.T) {
// 		router, _, _, _, tokenHandler := setupTestEnv()
// 		accessToken := "testaccesstoken123456789"
// 		refreshToken := "testrefreshtoken123456789"
// 		tokenHandler.On("CreateTokenPair", 0, "test").Return(&auth.TokenPair{
// 			AccessToken:  accessToken,
// 			RefreshToken: refreshToken,
// 		}, nil).Once()

// 		router.Use(AuthMiddleware(tokenHandler))

// 		tokenHandler.On("VerifyToken", accessToken).Return(&auth.TokenClaims{
// 			Username: "test",
// 			UserID:   0,
// 		}, nil).Once()

// 		router.GET("/protected", func(c *gin.Context) {
// 			userID := c.GetInt("userID")
// 			username := c.GetString("username")
// 			c.JSON(200, gin.H{
// 				"msg":      "authorized",
// 				"userID":   userID,
// 				"username": username,
// 			})
// 		})

// 		req := newRequest(t, "GET", "/protected", nil)
// 		req.Header.Set("Authorization", "Bearer "+accessToken)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.JSONEq(t, `{
// 			"msg": "authorized",
// 			"userID": 0,
// 			"username": "test"
// 		}`, w.Body.String())

// 	})

// 	t.Run("Test_AuthMiddleware_With_No_Authorization_Header", func(t *testing.T) {
// 		router, _, _, _, tokenHandler := setupTestEnv()

// 		router.Use(AuthMiddleware(tokenHandler))

// 		router.GET("/protected", func(c *gin.Context) {
// 			userID := c.GetInt("userID")
// 			username := c.GetString("username")
// 			c.JSON(200, gin.H{
// 				"msg":      "authorized",
// 				"userID":   userID,
// 				"username": username,
// 			})
// 		})

// 		req := newRequest(t, "GET", "/protected", nil)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusUnauthorized, w.Code)

// 	})

// 	t.Run("Test_AuthMiddleware_With_Invalid_Token", func(t *testing.T) {
// 		router, _, _, _, tokenHandler := setupTestEnv()
// 		accessToken := "testaccesstoken123456789"
// 		refreshToken := "testrefreshtoken123456789"
// 		tokenHandler.On("CreateTokenPair", 0, "test").Return(&auth.TokenPair{
// 			AccessToken:  accessToken,
// 			RefreshToken: refreshToken,
// 		}, nil).Once()

// 		router.Use(AuthMiddleware(tokenHandler))

// 		tokenHandler.On("VerifyToken", accessToken).Return(&auth.TokenClaims{}, errors.New("test error")).Once()

// 		router.GET("/protected", func(c *gin.Context) {
// 			userID := c.GetInt("userID")
// 			username := c.GetString("username")
// 			c.JSON(200, gin.H{
// 				"msg":      "authorized",
// 				"userID":   userID,
// 				"username": username,
// 			})
// 		})

// 		req := newRequest(t, "GET", "/protected", nil)
// 		req.Header.Set("Authorization", "Bearer "+accessToken)
// 		w := applyRequest(router, req)

// 		assert.Equal(t, http.StatusUnauthorized, w.Code)

// 	})

// }
