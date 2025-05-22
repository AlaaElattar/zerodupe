package client

import (
	"fmt"
	"sync"
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

	var wg sync.WaitGroup
	errChan := make(chan error, len(chunks))
	semaphore := make(chan struct{}, 5)

	for i, chunk := range chunks {

		wg.Add(1)
		chunkIndex := i

		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

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
				errChan <- fmt.Errorf("failed to upload chunk %d: %w", chunkIndex+1, err)
				return
			}

			if response.HashMismatch {
				fmt.Printf("WARNING: Server detected hash mismatch for chunk %s\n", chunk.ChunkHash)
			}

		}()

	}

	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}

	return nil
}
