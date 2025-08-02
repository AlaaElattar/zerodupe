package model

// FileMetadata represents metadata for a complete file
type FileMetadata struct {
	ID       uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	FileHash string          `gorm:"index;unique;not null" json:"file_hash"`
	Chunks   []ChunkMetadata `gorm:"foreignKey:FileMetadataID;constraint:OnDelete:CASCADE" json:"chunks"`
}

// ChunkMetadata represents metadata for a single chunk
type ChunkMetadata struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FileMetadataID uint   `gorm:"index" json:"file_metadata_id"` // foreign key
	ChunkOrder     int    `json:"chunk_order"`
	ChunkHash      string `gorm:"uniqueIndex:idx_file_chunk_order,priority:2" json:"chunk_hash"`
}
