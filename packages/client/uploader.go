package client

import (
	"fmt"
	"zerodupe/packages/hasher"
)

// Uploader handles uploading files to the server
type Uploader struct {
	api API
}

// NewUploader creates a new uploader
func NewUploader(api API) *Uploader {
	return &Uploader{
		api: api,
	}
}

// UploadChunks uploads all chunks to the server
func (u *Uploader) UploadChunks(chunks []hasher.FileChunk, fileHash string, filePath string, existingChunks map[string]bool) error {
	if existingChunks == nil {
		return nil
	}

	for i, chunk := range chunks {
		var content []byte

		if existingChunks[chunk.ChunkHash] {
			fmt.Printf("Chunk %d/%d (Hash: %s) already exists. Sending metadata only.\n",
				i+1, len(chunks), chunk.ChunkHash)
			content = nil
		} else {
			fmt.Printf("Uploading chunk %d/%d (Order: %d, Size: %d bytes, Hash: %s)\n",
				i+1, len(chunks), chunk.ChunkOrder, len(chunk.Data), chunk.ChunkHash)
			content = chunk.Data
		}

		request := UploadRequest{
			FileHash:   fileHash,
			ChunkHash:  chunk.ChunkHash,
			ChunkOrder: chunk.ChunkOrder,
			Content:    content,
		}

		response, err := u.api.UploadChunk(request)
		if err != nil {
			return fmt.Errorf("failed to upload chunk %d: %w", i+1, err)
		}

		if response.HashMismatch {
			fmt.Printf("WARNING: Server detected hash mismatch for chunk %s\n", chunk.ChunkHash)
		}
	}

	return nil
}
