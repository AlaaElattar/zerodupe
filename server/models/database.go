package models

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB struct holds the database connection
type DB struct {
	db    *gorm.DB
	mutex sync.Mutex
}

// FileStatus struct holds the status of a file
type FileStatus struct {
	HashExists bool
	NameExists bool
}

// Connect connects to the database file
func Connect(file string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&File{})
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// SaveFile saves a file to the database
func (d *DB) SaveFile(filename string, fileHash string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	fileStatus := d.CheckFileStatus(filename, fileHash)

	// If both hash and name exist, check if chunk exists
	if fileStatus.HashExists && fileStatus.NameExists {
		return nil
	}

	// If hash exists but name doesn't, add a new file entry with this name
	if fileStatus.HashExists {
		file := File{
			FileHash: fileHash,
			Name:     filename,
		}
		if err := d.db.Create(&file).Error; err != nil {
			return err
		}
		return nil
	}

	tx := d.db.Begin()
	defer tx.Rollback()

	file := File{
		FileHash: fileHash,
		Name:     filename,
	}
	if err := tx.Create(&file).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

// CheckFileStatus checks if a file exists in the database
func (d *DB) CheckFileStatus(name string, hash string) FileStatus {
	status := FileStatus{}
	var file File
	hashResult := d.db.First(&file, "file_hash = ?", hash)
	status.HashExists = hashResult.Error == nil
	if status.HashExists {
		nameResult := d.db.First(&file, "file_hash = ? AND name = ?", hash, name)
		status.NameExists = nameResult.Error == nil
	} else {
		status.NameExists = false
	}
	return status
}

// GetFileByName returns a file by name
func (d *DB) GetFileHashByName(name string) (string, error) {
	var file File
	result := d.db.First(&file, "name = ?", name)
	if result.Error != nil {
		return "", result.Error
	}
	return file.FileHash, nil
}
