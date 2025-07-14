package model

// ChunkMetadata represents metadata for a single chunk
type ChunkMetadata struct {
	ChunkOrder int    `json:"chunk_order"`
	ChunkHash  string `json:"chunk_hash"`
}

// FileMetadata represents metadata for a complete file
type FileMetadata struct {
	Chunks []ChunkMetadata `json:"chunks"`
}
