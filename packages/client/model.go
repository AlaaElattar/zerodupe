package client

type UploadRequest struct {
	FileHash   string `json:"filehash" binding:"required"`
	ChunkHash  string `json:"chunkhash" binding:"required"`
	ChunkOrder int    `json:"chunk_order" binding:"required"`
	Content    []byte `json:"content"`
}

// UploadResponse represents a response from the server after uploading a chunk
type UploadResponse struct {
	Message      string `json:"message"`
	FileHash     string `json:"fileHash"`
	HashMismatch bool   `json:"hashMismatch"`
}

// CheckFileResponse represents a response from the server when checking if a file exists
type CheckFileResponse struct {
	Exists bool   `json:"exists"`
	Hash   string `json:"hash"`
}

// CheckChunksRequest represents a request to check if chunks exist on the server
type CheckChunksRequest struct {
	Hashes []string `json:"hashes"`
}

// CheckChunksResponse represents a response from the server when checking if chunks exist
type CheckChunksResponse struct {
	Exists  []string `json:"exists"`
	Missing []string `json:"missing"`
}
