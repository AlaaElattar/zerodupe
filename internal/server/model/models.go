package model

// UploadRequest represents a file upload request
type UploadRequest struct {
	FileHash   string `json:"file_hash" binding:"required"`
	ChunkHash  string `json:"chunk_hash" binding:"required"`
	ChunkOrder int    `json:"chunk_order" binding:"required"`
	Content    []byte `json:"content"`
}

// ChunkMetadata represents metadata for a single chunk
type ChunkMetadata struct {
	ChunkOrder int    `json:"chunk_order"`
	ChunkHash  string `json:"chunk_hash"`
}

// FileMetadata represents metadata for a complete file
type FileMetadata struct {
	Chunks []ChunkMetadata `json:"chunks"`
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

// AuthRequest represents a request for authentication (login/signup)
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
