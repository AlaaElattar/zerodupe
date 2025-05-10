package models

// File represents a file
type File struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	FileHash string `json:"file_hash" binding:"required"`
	Name     string `json:"file_name" binding:"required"`
}
