package model

// UploadRequest represents a file upload request
type UploadRequest struct {
	FileHash   string `json:"filehash" binding:"required"`
	ChunkHash  string `json:"chunkhash" binding:"required"`
	ChunkOrder int    `json:"chunk_order" binding:"required"`
	Content    []byte `json:"content"`
}

// ChunkMetadataRequest represents metadata for a single chunk
type ChunkMetadataRequest struct {
	ChunkOrder int    `json:"chunk_order"`
	ChunkHash  string `json:"chunk_hash"`
}

// FileMetadata represents metadata for a complete file
type FileMetadata struct {
	Chunks []ChunkMetadataRequest `json:"chunks"`
}

// UploadResponse represents a response to an upload request
type UploadResponse struct {
	Message      string `json:"message"`
	FileHash     string `json:"fileHash"`
	HashMismatch bool   `json:"hashMismatch"`
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

type DownloadFileResponse struct {
	FileHash    string   `json:"filehash" binding:"required"`
	ChunkHashes []string `json:"chunk_hashes"`
	ChunksCount int      `json:"chunks_count"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
