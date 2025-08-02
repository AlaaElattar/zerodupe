package filesystem

import (
	"os"
	"path/filepath"
	"testing"
	"zerodupe/pkg/hasher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupFileSystemStorage creates a temporary filesystem storage for testing
func setupFileSystemStorage(t *testing.T) (*FilesystemStorage, string) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "zerodupe-test")
	require.NoError(t, err, "Failed to create temporary directory")
	storage, err := NewFilesystemStorage(tempDir)
	require.NoError(t, err, "Failed to create filesystem storage")

	return storage, tempDir
}

// teardownFileSystemStorage removes the temporary filesystem storage
func teardownFileSystemStorage(t *testing.T, tempDir string) {
	t.Helper()

	err := os.RemoveAll(tempDir)
	require.NoError(t, err, "Failed to remove temporary directory")
}

func TestNewFilesystemStorage(t *testing.T) {
	t.Run("Test NewFilesystemStorage creates a new filesystem storage", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		require.NotNil(t, storage)
		require.Equal(t, tempDir, storage.storageDir)
		require.DirExists(t, tempDir)
		require.DirExists(t, filepath.Join(tempDir, "blocks"))
		// Check directory permissions
		for _, dir := range []string{tempDir, filepath.Join(tempDir, "blocks")} {
			info, err := os.Stat(dir)
			require.NoError(t, err)
			require.True(t, info.IsDir())
			require.True(t, info.Mode().Perm()&0700 == 0700, "Directory should be readable/writable.")
		}
	})

	t.Run("Test NewFilesystemStorage returns error for invalid directory", func(t *testing.T) {
		_, err := NewFilesystemStorage("/invalid")
		require.Error(t, err)
	})

	t.Run("Test NewFilesystemStorage returns error for read-only directory", func(t *testing.T) {
		_, err := NewFilesystemStorage("/root")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only file system")
	})
}

func TestCheckFileExists(t *testing.T) {

	t.Run("Test CheckFileExists for existing file in blocks directory", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		fileHash := "abcdtest567890"
		blockDir := filepath.Join(tempDir, "blocks", fileHash[:4])
		err := os.MkdirAll(blockDir, 0755)
		require.NoError(t, err)

		blockPath := filepath.Join(blockDir, fileHash)
		err = os.WriteFile(blockPath, []byte("test"), 0644)
		require.NoError(t, err)

		exists, err := storage.CheckFileExists(fileHash)
		require.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, fileHash, fileHash)

	})

	t.Run("Test CheckFileExists for non-existing file in blocks directory", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		fileHash := "abcdtest567890"
		exists, err := storage.CheckFileExists(fileHash)
		require.NoError(t, err)
		assert.False(t, exists)

	})

}

func TestCheckChunkExists(t *testing.T) {

	t.Run("Test CheckChunkExists for existing chunk", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		chunkHash := "abcdtest567890"
		blockDir := filepath.Join(tempDir, "blocks", chunkHash[:4])
		err := os.MkdirAll(blockDir, 0755)
		require.NoError(t, err)

		blockPath := filepath.Join(blockDir, chunkHash)
		err = os.WriteFile(blockPath, []byte("test"), 0644)
		require.NoError(t, err)

		existing, missing, err := storage.CheckChunkExists([]string{chunkHash})
		require.NoError(t, err)
		assert.Equal(t, 1, len(existing))
		assert.Equal(t, 0, len(missing))

	})

	t.Run("Test CheckChunkExists for non-existing chunk", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		chunkHash := "abcdtest567890"
		existing, missing, err := storage.CheckChunkExists([]string{chunkHash})
		require.NoError(t, err)
		assert.Equal(t, 0, len(existing))
		assert.Equal(t, 1, len(missing))
	})

	t.Run("Test CheckChunkExists for multiple chunks", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		chunkHash1 := "abcdtest567890"
		chunkHash2 := "efghijklmnop"
		chunkHash3 := "qrstuvwxyz12"

		blockDir1 := filepath.Join(tempDir, "blocks", chunkHash1[:4])
		err := os.MkdirAll(blockDir1, 0755)
		require.NoError(t, err)

		blockPath1 := filepath.Join(blockDir1, chunkHash1)
		err = os.WriteFile(blockPath1, []byte("test"), 0644)
		require.NoError(t, err)

		blockDir2 := filepath.Join(tempDir, "blocks", chunkHash2[:4])
		err = os.MkdirAll(blockDir2, 0755)
		require.NoError(t, err)

		blockPath2 := filepath.Join(blockDir2, chunkHash2)
		err = os.WriteFile(blockPath2, []byte("test"), 0644)
		require.NoError(t, err)

		existing, missing, err := storage.CheckChunkExists([]string{chunkHash1, chunkHash2, chunkHash3})
		require.NoError(t, err)
		assert.Equal(t, 2, len(existing))
		assert.Equal(t, 1, len(missing))
	})

}

func TestSaveChunkData(t *testing.T) {

	t.Run("Test SaveChunkData for new chunk", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		content := []byte("test")
		chunkHash := hasher.CalculateChunkHash(content)

		calculatedHash, err := storage.SaveChunkData(chunkHash, content)

		require.NoError(t, err)
		assert.Equal(t, hasher.CalculateChunkHash(content), calculatedHash)

		blockPath := filepath.Join(tempDir, "blocks", chunkHash[:4], chunkHash)
		assert.FileExists(t, blockPath)

		content, err = os.ReadFile(blockPath)
		require.NoError(t, err)
		assert.Equal(t, content, []byte("test"))
	})

	t.Run("Test SaveChunkData for existing chunk", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		content := []byte("test")
		chunkHash := hasher.CalculateChunkHash(content)

		calculatedHash, err := storage.SaveChunkData(chunkHash, content)
		require.NoError(t, err)
		assert.Equal(t, hasher.CalculateChunkHash(content), calculatedHash)

		calculatedHash, err = storage.SaveChunkData(chunkHash, content)
		require.NoError(t, err)
		assert.Equal(t, hasher.CalculateChunkHash(content), calculatedHash)

		blockPath := filepath.Join(tempDir, "blocks", chunkHash[:4], chunkHash)
		assert.FileExists(t, blockPath)

		content, err = os.ReadFile(blockPath)
		require.NoError(t, err)
		assert.Equal(t, content, []byte("test"))
	})

	t.Run("Test SaveChunkData with mismatched hash", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		content := []byte("test")
		incorrectHash := "incorrectHash"

		calculatedHash, err := storage.SaveChunkData(incorrectHash, content)

		require.NoError(t, err)
		assert.NotEqual(t, incorrectHash, calculatedHash)

		blockPath := filepath.Join(tempDir, "blocks", incorrectHash[:4], incorrectHash)
		assert.FileExists(t, blockPath)

		content, err = os.ReadFile(blockPath)
		require.NoError(t, err)
		assert.Equal(t, content, []byte("test"))

	})
}

func TestGetChunkData(t *testing.T) {
	t.Run("Test GetChunkData for existing chunk", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		content := []byte("test")
		chunkHash := hasher.CalculateChunkHash(content)

		calculatedHash, err := storage.SaveChunkData(chunkHash, content)
		require.NoError(t, err)
		assert.Equal(t, hasher.CalculateChunkHash(content), calculatedHash)

		retrievedContent, err := storage.GetChunkData(chunkHash)
		require.NoError(t, err)
		assert.Equal(t, content, retrievedContent)
	})

	t.Run("Test GetChunkData for non-existing chunk", func(t *testing.T) {
		storage, tempDir := setupFileSystemStorage(t)
		defer teardownFileSystemStorage(t, tempDir)

		chunkHash := "abcdtest567890"

		_, err := storage.GetChunkData(chunkHash)
		assert.Error(t, err)
		assert.Equal(t, "chunk not found: stat "+filepath.Join(tempDir, "blocks", chunkHash[:4], chunkHash)+": no such file or directory", err.Error())
	})
}
