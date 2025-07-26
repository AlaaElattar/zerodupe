package storage

import (
	"errors"
	"fmt"
	"zerodupe/internal/server/model"

	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func NewGormStorage(dialector gorm.Dialector) (DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate models
	err = db.AutoMigrate(&model.User{}, &model.FileMetadata{}, &model.ChunkMetadata{})
	if err != nil {
		return nil, err
	}

	return &GormDB{db: db}, nil
}

// CreateUser creates a new user
func (g *GormDB) CreateUser(user *model.User) error {
	return g.db.Create(user).Error

}

// GetUserByUsername gets a user by username
func (g *GormDB) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := g.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Close closes the database connection
func (g *GormDB) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (g *GormDB) SaveChunkMetadata(fileHash, chunkHash string, chunkOrder int) error {
	var fileMetadata model.FileMetadata

	err := g.db.Where("file_hash = ?", fileHash).First(&fileMetadata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fileMetadata = model.FileMetadata{FileHash: fileHash}
			if err := g.db.Create(&fileMetadata).Error; err != nil {
				return fmt.Errorf("failed to create file metadata: %w", err)

			}
		} else {
			return fmt.Errorf("failed to query file metadata: %w", err)

		}
	}

	var existingChunk model.ChunkMetadata
	err = g.db.Where("file_metadata_id = ? AND chunk_order = ? AND chunk_hash = ?", fileMetadata.ID, chunkOrder, chunkHash).
		First(&existingChunk).Error
	if err == nil {
		// Chunk already exists
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing chunk: %w", err)
	}

	newChunk := model.ChunkMetadata{
		FileMetadataID: fileMetadata.ID,
		ChunkOrder:     chunkOrder,
		ChunkHash:      chunkHash,
	}
	if err := g.db.Create(&newChunk).Error; err != nil {
		return fmt.Errorf("failed to save chunk metadata: %w", err)
	}

	return nil

}

func (g *GormDB) GetFileMetadata(fileHash string) (*model.FileMetadata, error) {
	var fileMetadata model.FileMetadata

	err := g.db.Preload("Chunks").Where("file_hash = ?", fileHash).First(&fileMetadata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	return &fileMetadata, nil
}

func (g *GormDB) CheckFileExists(fileHash string) (bool, error) {
	var file model.FileMetadata
	err := g.db.
		Select("id").
		Where("file_hash = ?", fileHash).
		First(&file).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("database error checking file existence: %w", err)
	}

	return true, nil
}
