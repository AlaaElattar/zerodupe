package client

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"zerodupe/pkg/hasher"
)

// ChunkDownloader handles downloading chunks from server
type ChunkDownloader struct {
	api API
}

// NewDownloader creates a new downloader
func NewDownloader(api API) *ChunkDownloader {
	return &ChunkDownloader{
		api: api,
	}
}

func (d *ChunkDownloader) DownloadChunks(hashes *DownloadFileHashesResponse) ([][]byte, error) {

	var wg sync.WaitGroup
	resultChan := make(chan ChunkDownloadResult, len(hashes.ChunkHashes))
	semaphore := make(chan struct{}, 5)

	var downloadContent atomic.Int32
	totalChunks := len(hashes.ChunkHashes)
	progressTicker := time.NewTicker(500 * time.Millisecond)
	defer progressTicker.Stop()

	go reportDownloadProgress(&downloadContent, totalChunks, progressTicker)

	for i, hash := range hashes.ChunkHashes {
		wg.Add(1)
		chunkIndex := i
		currentHash := hash
		go func(hash string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fmt.Printf("Downloading chunk %d/%d with hash: %s\n", chunkIndex+1, hashes.ChunksCount, hash)

			response, err := d.api.DownloadChunk(hash)

			if err == nil {
				calculatedHash := hasher.CalculateChunkHash(response)
				if calculatedHash != hash {
					fmt.Printf("WARNING: Hash mismatch for chunk %d. Expected: %s, Got: %s\n",
						chunkIndex+1, hash, calculatedHash)
				}
			}

			resultChan <- ChunkDownloadResult{
				index:   chunkIndex,
				content: response,
				err:     err,
			}

		}(currentHash)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make([][]byte, len(hashes.ChunkHashes))
	for result := range resultChan {
		if result.err != nil {
			return nil, fmt.Errorf("failed to download chunk %d: %w", result.index+1, result.err)
		}
		results[result.index] = result.content
	}

	if len(results) != hashes.ChunksCount {
		return nil, fmt.Errorf("downloaded chunks count does not match expected count")
	}

	return results, nil
}

// reportDownloadProgress displays download progress at regular intervals
func reportDownloadProgress(counter *atomic.Int32, total int, ticker *time.Ticker) {
	for range ticker.C {
		current := counter.Load()
		if current >= int32(total) {
			break
		}
		percentage := float64(current) / float64(total) * 100
		fmt.Printf("\rDownloading: %d/%d chunks (%.1f%%)", current, total, percentage)
	}
}