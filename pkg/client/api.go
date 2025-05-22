package client

type API interface {
	// CheckFileExists checks if a file exists on the server
	CheckFileExists(fileHash string) (bool, error)

	// CheckChunksExists checks if chunks exist on the server
	CheckChunksExists(hashes []string) ([]string, []string, error)

	// UploadChunk uploads a chunk to the server
	UploadChunk(request UploadRequest) (*UploadResponse, error)

	// DownloadFile downloads a file from the server
	DownloadFileHashes(fileHash string) (*DownloadFileHashesResponse, error)

	DownloadChunkContent(chunkHash string) ([]byte, error)
}
