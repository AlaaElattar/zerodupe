package client

// ChunkUploadRequest represents a request to upload a chunk to the server
type ChunkUploadRequest struct {
	FileHash   string `json:"filehash" binding:"required"`
	ChunkHash  string `json:"chunkhash" binding:"required"`
	ChunkOrder int    `json:"chunk_order" binding:"required"`
	Content    []byte `json:"content"`
}

// ChunkUploadResponse represents a response from the server after uploading a chunk
type ChunkUploadResponse struct {
	Message      string `json:"message"`
	FileHash     string `json:"fileHash"`
	HashMismatch bool   `json:"hashMismatch"`
}

// FileExistsResponse represents a response from the server when checking if a file exists
type FileExistsResponse struct {
	Exists bool   `json:"exists"`
	Hash   string `json:"hash"`
}

// MissingChunksRequest represents a request to check if chunks exist on the server
type MissingChunksRequest struct {
	Hashes []string `json:"hashes"`
}

// MissingChunksResponse represents a response from the server when checking if chunks exist
type MissingChunksResponse struct {
	Exists  []string `json:"exists"`
	Missing []string `json:"missing"`
}

type DownloadFileHashesResponse struct {
	FileHash    string   `json:"filehash" binding:"required"`
	ChunkHashes []string `json:chunk_hashes`
	ChunksCount int      `json:chunks_count`
}

// ChunkDownloadResult represents the result of downloading a chunk
type ChunkDownloadResult struct {
	index   int
	content []byte
	err     error
}
