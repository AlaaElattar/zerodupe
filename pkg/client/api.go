package client

type API interface {
	// CheckFileExists checks if a file exists on the server
	CheckFileExists(fileHash string) (bool, error)

	// GetMissingChunks checks if chunks exist on the server
	GetMissingChunks(hashes []string) ([]string, error)

	// UploadChunk uploads a chunk to the server
	UploadChunk(request ChunkUploadRequest) (*ChunkUploadResponse, error)

	// DownloadFile downloads a file from the server
	GetFileChunks(fileHash string) (*DownloadFileHashesResponse, error)

	DownloadChunk(chunkHash string) ([]byte, error)
}
