package client

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"zerodupe/pkg/hasher"
)

// ChunkUploader handles uploading chunks to the server
type ChunkUploader struct {
	api API
}

// NewUploader creates a new uploader
func NewUploader(api API) *ChunkUploader {
	return &ChunkUploader{
		api: api,
	}
}

// UploadChunks uploads all chunks to the server
func (u *ChunkUploader) UploadChunks(chunks []hasher.FileChunk, fileHash string, filePath string, existingChunks map[string]bool) error {
	if existingChunks == nil {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(chunks))
	semaphore := make(chan struct{}, 5)

	var uploadedCount atomic.Int32
	totalChunks := len(chunks)
	progressTicker := time.NewTicker(500 * time.Millisecond)
	defer progressTicker.Stop()

	go reportUploadProgress(&uploadedCount, totalChunks, progressTicker)

	for i, chunk := range chunks {

		wg.Add(1)
		chunkIndex := i
		currentChunk := chunk

		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			var content []byte

			if existingChunks[currentChunk.ChunkHash] {
				fmt.Printf("Chunk %d/%d (Hash: %s) already exists. Sending metadata only.\n",
					i+1, len(chunks), currentChunk.ChunkHash)
				content = nil
			} else {
				fmt.Printf("Uploading chunk %d/%d (Order: %d, Size: %d bytes, Hash: %s)\n",
					i+1, len(chunks), currentChunk.ChunkOrder, len(currentChunk.Data), currentChunk.ChunkHash)
				content = currentChunk.Data
			}

			request := ChunkUploadRequest{
				FileHash:   fileHash,
				ChunkHash:  currentChunk.ChunkHash,
				ChunkOrder: currentChunk.ChunkOrder,
				Content:    content,
			}

			response, err := u.api.UploadChunk(request)
			if err != nil {
				errChan <- fmt.Errorf("failed to upload chunk %d: %w", chunkIndex+1, err)
				return
			}

			if response.HashMismatch {
				fmt.Printf("WARNING: Server detected hash mismatch for chunk %s\n", currentChunk.ChunkHash)
			}

		}()

	}

	wg.Wait()
	close(errChan)
	fmt.Printf("\rUpload complete: %d/%d chunks (100%%)      \n", totalChunks, totalChunks)

	for err := range errChan {
		return err
	}

	return nil
}

// reportUploadProgress displays upload progress at regular intervals
func reportUploadProgress(counter *atomic.Int32, total int, ticker *time.Ticker) {
	for range ticker.C {
		current := counter.Load()
		if current >= int32(total) {
			break
		}
		percentage := float64(current) / float64(total) * 100
		fmt.Printf("\rUploading: %d/%d chunks (%.1f%%)", current, total, percentage)
	}
}
