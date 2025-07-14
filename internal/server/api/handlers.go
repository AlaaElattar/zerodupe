package api

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zerodupe/internal/server/auth"
	"zerodupe/internal/server/model"
	"zerodupe/internal/server/storage"
)

// Handler handles all API requests
type Handler struct {
	fileStorage  storage.FileStorage
	userStorage  storage.UserStorage
	tokenHandler auth.TokenManager
}

func NewHandler(fileStorage storage.FileStorage, userStorage storage.UserStorage, tokenHandler auth.TokenManager) *Handler {
	return &Handler{
		fileStorage:  fileStorage,
		userStorage:  userStorage,
		tokenHandler: tokenHandler,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// SignUpRequest holds data needed for signup
type SignUpRequest struct {
	Username        string `json:"username" binding:"required" example:"john_doe"`
	Password        string `json:"password" binding:"required" example:"password123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" example:"password123"`
}

// RefreshTokenRequest represents the request body for refreshing an access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UploadRequest represents a file upload request
type UploadRequest struct {
	FileHash   string `json:"file_hash" binding:"required"`
	ChunkHash  string `json:"chunk_hash" binding:"required"`
	ChunkOrder int    `json:"chunk_order" binding:"required"`
	Content    []byte `json:"content"`
}

// UploadResponse represents a response to an upload request
type UploadResponse struct {
	Message      string `json:"message"`
	FileHash     string `json:"file_hash"`
	HashMismatch bool   `json:"hash_mismatch"`
}

// CheckFileResponse represents a response to a file existence check
type CheckFileResponse struct {
	Exists bool   `json:"exists"`
	Hash   string `json:"hash"`
}

// CheckChunksResponse represents a response to a chunks existence check
type CheckChunksResponse struct {
	Missing []string `json:"missing"`
}

// DownloadFileResponse represents a response to a file download request
type DownloadFileResponse struct {
	FileHash    string   `json:"file_hash" binding:"required"`
	ChunkHashes []string `json:"chunk_hashes"`
	ChunksCount int      `json:"chunks_count"`
}

// CheckChunksRequest represents the request body for checking chunk hashes
type CheckChunksRequest struct {
	Hashes []string `json:"hashes" binding:"required"`
}

// @Summary Register a new user
// @Description Create a new user account with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "User registration data"
// @Success 200 {string} string "user registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or password mismatch"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/signup [post]
func (h *Handler) SignUpHandler(c *gin.Context) {
	var request SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password and confirm password don't match"})
		return
	}

	// Check if username already exists
	_, err := h.userStorage.GetUserByUsername(request.Username)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	password, err := auth.HashAndSaltPassword([]byte(request.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	err = h.userStorage.CreateUser(&model.User{
		Username: request.Username,
		Password: password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, "user registered successfully")

}

// @Summary Login user
// @Description Authenticate user and return access tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} auth.TokenPair "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 404 {object} map[string]interface{} "User does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.userStorage.GetUserByUsername(request.Username)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	// Generate tokens
	tokenPair, err := h.tokenHandler.CreateTokenPair(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// @Summary Refresh access token
// @Description Refresh an expired access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]string "New access token"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Router /auth/refresh [post]
func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	var request RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	accessToken, err := h.tokenHandler.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// @Summary Upload file chunk
// @Description Upload a file chunk for deduplication storage
// @Tags files
// @Accept json
// @Produce json
// @Param request body UploadRequest true "File chunk data"
// @Success 200 {object} UploadResponse "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 404 {object} map[string]interface{} "Chunk does not exist"
// @Failure 500 {object} map[string]interface{} "Failed to save chunk data"
// @Router /upload [post]
func (h *Handler) UploadFileHandler(c *gin.Context) {
	var request UploadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// special case: one chunk file, save chunk directly without metadata
	if request.FileHash == request.ChunkHash {
		_, err := h.fileStorage.SaveChunkData(request.ChunkHash, request.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk data"})
			return
		}
		response := UploadResponse{
			Message:      "File uploaded successfully",
			FileHash:     request.FileHash,
			HashMismatch: false,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	log.Printf("Received chunk: Hash=%s, ChunkOrder=%d\n",
		request.FileHash, request.ChunkOrder)

	if err := h.fileStorage.SaveChunkMetadata(request.FileHash, request.ChunkHash, request.ChunkOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk metadata"})
		return
	}

	hashMismatch := false
	if (len(request.Content)) > 0 {
		calculatedHash, err := h.fileStorage.SaveChunkData(request.ChunkHash, request.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk data"})
			return
		}
		hashMismatch = calculatedHash != request.ChunkHash
	}
	if (len(request.Content)) == 0 {
		exists, _, err := h.fileStorage.CheckChunkExists([]string{request.ChunkHash})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chunk existence"})
			return
		}
		if len(exists) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chunk does not exist"})
			return
		}
	}

	response := UploadResponse{
		Message:      "File uploaded successfully",
		FileHash:     request.FileHash,
		HashMismatch: hashMismatch,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Check if file exists
// @Description Check if a file with the given hash exists on the server
// @Tags files
// @Accept json
// @Produce json
// @Param filehash path string true "File hash" minlength(4)
// @Success 200 {object} CheckFileResponse "File existence status"
// @Failure 400 {object} map[string]interface{} "Invalid file hash"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /check/{filehash} [get]
func (h *Handler) CheckFileHashHandler(c *gin.Context) {
	fileHash := c.Param("filehash")

	if len(fileHash) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file hash"})
		return
	}

	exists, err := h.fileStorage.CheckFileExists(fileHash)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := CheckFileResponse{
		Exists: exists,
		Hash:   fileHash,
	}

	c.JSON(200, response)
}

// @Summary Check chunk existence
// @Description Check which chunks exist and return missing ones
// @Tags files
// @Accept json
// @Produce json
// @Param request body CheckChunksRequest true "Chunk hashes to check"
// @Success 200 {object} CheckChunksResponse "Missing chunks"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /check [post]
func (h *Handler) CheckChunkHashesHandler(c *gin.Context) {
	var request CheckChunksRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, missing, err := h.fileStorage.CheckChunkExists(request.Hashes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := CheckChunksResponse{
		Missing: missing,
	}

	c.JSON(200, response)
}

// @Summary Download file metadata
// @Description Get file metadata including ordered chunk hashes for download
// @Tags files
// @Accept json
// @Produce json
// @Param hash path string true "File hash"
// @Success 200 {object} DownloadFileResponse "File metadata"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /download/{hash} [get]
func (h *Handler) DownloadFileHandler(c *gin.Context) {
	fileHash := c.Param("hash")

	log.Println("Downloading file with hash: " + fileHash)
	metadata, err := h.fileStorage.GetFileMetadata(fileHash)

	if err != nil {
		// If file metadata not found, check if it's a single chunk file
		_, err := h.fileStorage.GetChunkData(fileHash)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		result := DownloadFileResponse{
			FileHash:    fileHash,
			ChunkHashes: []string{fileHash},
			ChunksCount: 1,
		}
		c.JSON(http.StatusOK, result)

		return
	}

	sort.Slice(metadata.Chunks, func(i, j int) bool {
		return metadata.Chunks[i].ChunkOrder < metadata.Chunks[j].ChunkOrder
	})

	var orderedHashes []string
	for _, chunk := range metadata.Chunks {
		orderedHashes = append(orderedHashes, chunk.ChunkHash)
	}

	// return chunks hashes ordered
	result := DownloadFileResponse{
		FileHash:    fileHash,
		ChunkHashes: orderedHashes,
		ChunksCount: len(orderedHashes),
	}

	c.JSON(http.StatusOK, result)

}

// @Summary Get chunk content
// @Description Download the content of a specific chunk
// @Tags files
// @Accept json
// @Produce octet-stream
// @Param hash path string true "Chunk hash"
// @Success 200 {file} binary "Chunk content"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /chunk/{hash} [get]
func (h *Handler) GetChunkContent(c *gin.Context) {
	chunkHash := c.Param("hash")
	content, err := h.fileStorage.GetChunkData(chunkHash)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", content)
}
