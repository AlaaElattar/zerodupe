package client

import (
	"fmt"
	"sync"
	"zerodupe/packages/hasher"
)

// Downloader handles downloading chunks from server
type Downloader struct {
	api API
}

// NewDownloader creates a new downloader
func NewDownloader(api API) *Downloader {
	return &Downloader{
		api: api,
	}
}

func (d *Downloader) DownloadChunks(hashes *DownloadFileHashesResponse) ([][]byte, error) {

	var wg sync.WaitGroup
	resultChan := make(chan chunkResult, len(hashes.ChunkHashes))
	semaphore := make(chan struct{}, 5)

	for i, hash := range hashes.ChunkHashes {
		wg.Add(1)
		chunkIndex := i
		chunkHash := hash
		go func(hash string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fmt.Printf("Downloading chunk %d/%d with hash: %s\n", chunkIndex+1, hashes.ChunksCount, hash)

			response, err := d.api.DownloadChunkContent(hash)

			if err == nil {
				calculatedHash := hasher.CalculateChunkHash(response)
				if calculatedHash != hash {
					fmt.Printf("WARNING: Hash mismatch for chunk %d. Expected: %s, Got: %s\n",
						chunkIndex+1, hash, calculatedHash)
				}
			}

			resultChan <- chunkResult{
				index:   chunkIndex,
				content: response,
				err:     err,
			}

		}(chunkHash)
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
