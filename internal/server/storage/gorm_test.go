package storage

import (
	"testing"
	"zerodupe/internal/server/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestGormDB(t *testing.T) *GormDB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.User{}, &model.FileMetadata{}, &model.ChunkMetadata{})
	require.NoError(t, err)

	return &GormDB{db: db}
}

func TestCreateUser(t *testing.T) {
	t.Run("Test CreateUser creates a new user", func(t *testing.T) {
		db := setupTestGormDB(t)
		user := &model.User{
			Username: "test",
			Password: []byte("test"),
		}
		err := db.CreateUser(user)

		assert.NoError(t, err)
		assert.Equal(t, "test", user.Username)
		assert.NotEmpty(t, user.Password)

		retrievedUser, err := db.GetUserByUsername("test")
		assert.NoError(t, err)
		assert.Equal(t, user.Username, retrievedUser.Username)
		assert.Equal(t, user.Password, retrievedUser.Password)
	})

	t.Run("Test CreateUser returns error for existing user", func(t *testing.T) {
		db := setupTestGormDB(t)

		user1 := &model.User{
			Username: "test",
			Password: []byte("test"),
		}
		err := db.CreateUser(user1)
		assert.NoError(t, err)
		assert.Equal(t, "test", user1.Username)
		assert.NotEmpty(t, user1.Password)
		assert.NotEqual(t, user1.Password, "test")

		user2 := &model.User{
			Username: "test",
			Password: []byte("test"),
		}
		err = db.CreateUser(user2)
		assert.Error(t, err)
	})
}
func TestGetUserByUsername(t *testing.T) {
	t.Run("Test GetUserByUsername returns user for existing user", func(t *testing.T) {
		db := setupTestGormDB(t)
		user := &model.User{
			Username: "test",
			Password: []byte("test"),
		}
		err := db.CreateUser(user)
		assert.NoError(t, err)

		retrievedUser, err := db.GetUserByUsername("test")
		assert.NoError(t, err)
		assert.Equal(t, user.Username, retrievedUser.Username)
		assert.Equal(t, user.Password, retrievedUser.Password)
	})

	t.Run("Test GetUserByUsername returns error for non-existing user", func(t *testing.T) {
		db := setupTestGormDB(t)
		_, err := db.GetUserByUsername("test")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestSaveChunkMetadata(t *testing.T) {

	t.Run("Test SaveChunkMetadata with new file and chunk", func(t *testing.T) {
		db := setupTestGormDB(t)
		fileHash := "hashabcd12345678910"
		chunkHash := "chunkhash12345678910"
		chunkOrder := 1

		err := db.SaveChunkMetadata(fileHash, chunkHash, chunkOrder)
		assert.NoError(t, err)
		fileMeta, err := db.GetFileMetadata(fileHash)
		assert.NoError(t, err)
		assert.NotNil(t, fileMeta)
		assert.Equal(t, fileHash, fileMeta.FileHash)
		assert.Len(t, fileMeta.Chunks, 1)
		assert.Equal(t, chunkHash, fileMeta.Chunks[0].ChunkHash)
		assert.Equal(t, chunkOrder, fileMeta.Chunks[0].ChunkOrder)

	})
	t.Run("Test SaveChunkMetadata adds another chunk to same file", func(t *testing.T) {
		db := setupTestGormDB(t)
		fileHash := "filehash1"
		_ = db.SaveChunkMetadata(fileHash, "chunk1", 0)

		err := db.SaveChunkMetadata(fileHash, "chunk2", 1)
		assert.NoError(t, err)

		fileMeta, err := db.GetFileMetadata(fileHash)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(fileMeta.Chunks))
	})

	t.Run("Test SaveChunkMetadata adds duplicated chunk", func(t *testing.T) {
		db := setupTestGormDB(t)
		fileHash := "filehash2"
		chunkHash := "chunk3"
		err := db.SaveChunkMetadata(fileHash, chunkHash, 2)
		assert.NoError(t, err)

		err = db.SaveChunkMetadata(fileHash, chunkHash, 2)
		assert.NoError(t, err)

		fileMeta, err := db.GetFileMetadata(fileHash)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(fileMeta.Chunks))
	})

	t.Run("Test SaveChunkMetadata with different file", func(t *testing.T) {
		db := setupTestGormDB(t)
		err := db.SaveChunkMetadata("filehashA", "chunkA1", 0)
		assert.NoError(t, err)

		err = db.SaveChunkMetadata("filehashB", "chunkB1", 0)
		assert.NoError(t, err)

		meta1, err := db.GetFileMetadata("filehashA")
		assert.NoError(t, err)

		meta2, err := db.GetFileMetadata("filehashB")
		assert.NoError(t, err)

		assert.Equal(t, 1, len(meta1.Chunks))
		assert.Equal(t, 1, len(meta2.Chunks))
	})
}

func TestGetFileMetadata(t *testing.T) {

	t.Run("Test GetFileMetadata", func(t *testing.T) {
		db := setupTestGormDB(t)

		file := model.FileMetadata{
			FileHash: "abc123",
			Chunks: []model.ChunkMetadata{
				{ChunkOrder: 1, ChunkHash: "chunk1"},
				{ChunkOrder: 2, ChunkHash: "chunk2"},
			},
		}
		err := db.db.Create(&file).Error
		require.NoError(t, err)
		got, err := db.GetFileMetadata("abc123")

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "abc123", got.FileHash)
		assert.Len(t, got.Chunks, 2)
		assert.Equal(t, "chunk1", got.Chunks[0].ChunkHash)
		assert.Equal(t, 1, got.Chunks[0].ChunkOrder)
		assert.Equal(t, "chunk2", got.Chunks[1].ChunkHash)
		assert.Equal(t, 2, got.Chunks[1].ChunkOrder)

	})

	t.Run("Test GetFileMetadata for non-existing file", func(t *testing.T) {
		db := setupTestGormDB(t)

		got, err := db.GetFileMetadata("abc123")
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestCheckFileExists(t *testing.T) {
	t.Run("Test CheckFileExists", func(t *testing.T) {
		db := setupTestGormDB(t)

		fileHash := "filehash123"
		err := db.db.Create(&model.FileMetadata{FileHash: fileHash}).Error
		assert.NoError(t, err)

		exists, err := db.CheckFileExists(fileHash)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Test CheckFileExists for non-existing file", func(t *testing.T) {
		db := setupTestGormDB(t)

		exists, err := db.CheckFileExists("nonexistent_hash")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
