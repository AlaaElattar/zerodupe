package models

// Chunk represents a chunk of a file
type Chunk struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	FileHash   string `json:"file_hash" gorm:"index:idx_file_hash" binding:"required"`
	ChunkOrder int    `json:"chunk_order" gorm:"index:idx_chunk_order" binding:"required"`
	ChunkHash  string `json:"chunkhash" gorm:"unique" binding:"required"`
	Chunk      []byte `json:"chunk" binding:"required"`
}
