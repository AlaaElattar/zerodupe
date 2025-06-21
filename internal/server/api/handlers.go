package api

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

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

// SignUpHandler handles user signup requests
func (h *Handler) SignUpHandler(c *gin.Context) {
	var request model.AuthRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if username already exists
	_, err := h.userStorage.LoginUser(request.Username, request.Password)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Create user
	user := &model.User{
		Username: request.Username,
	}
	err = h.userStorage.CreateUser(user, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Login users and get tokens
	user, err = h.userStorage.LoginUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	// Generate tokens
	tokenPair, err := h.tokenHandler.CreateTokenPair(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// LoginHandler handles user login requests
func (h *Handler) LoginHandler(c *gin.Context) {
	var request model.AuthRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.userStorage.LoginUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	tokenPair, err := h.tokenHandler.CreateTokenPair(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// RefreshTokenHandler handles token refresh requests
func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
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

// UploadFileHandler handles file upload requests
func (h *Handler) UploadFileHandler(c *gin.Context) {
	var request model.UploadRequest
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
		response := model.UploadResponse{
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

	response := model.UploadResponse{
		Message:      "File uploaded successfully",
		FileHash:     request.FileHash,
		HashMismatch: hashMismatch,
	}

	c.JSON(http.StatusOK, response)
}

// CheckFileHashHandler checks if a file exists on the server
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

	response := model.CheckFileResponse{
		Exists: exists,
		Hash:   fileHash,
	}

	c.JSON(200, response)
}

// CheckChunkHashesHandler checks if chunks exist on the server and returns the missing ones
func (h *Handler) CheckChunkHashesHandler(c *gin.Context) {
	var request struct {
		Hashes []string `json:"hashes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, missing, err := h.fileStorage.CheckChunkExists(request.Hashes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := model.CheckChunksResponse{
		Missing: missing,
	}

	c.JSON(200, response)
}

// DownloadFileHandler returns the chunks hashes ordered
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

		result := model.DownloadFileResponse{
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
	result := model.DownloadFileResponse{
		FileHash:    fileHash,
		ChunkHashes: orderedHashes,
		ChunksCount: len(orderedHashes),
	}

	c.JSON(http.StatusOK, result)

}

// GetChunkContent returns the content of a chunk
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
