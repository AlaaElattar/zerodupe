package api

import (
	"zerodupe/internal/server/auth"
	"zerodupe/internal/server/model"

	"github.com/stretchr/testify/mock"
)

type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) CheckFileExists(fileHash string) (bool, error) {
	args := m.Called(fileHash)
	return args.Bool(0), args.Error(1)
}

func (m *MockFileStorage) CheckChunkExists(hashes []string) ([]string, []string, error) {
	args := m.Called(hashes)
	return args.Get(0).([]string), args.Get(1).([]string), args.Error(2)
}

func (m *MockFileStorage) SaveChunkMetadata(fileHash, chunkHash string, chunkOrder int) error {
	args := m.Called(fileHash, chunkHash, chunkOrder)
	return args.Error(0)
}

func (m *MockFileStorage) SaveChunkData(chunkHash string, content []byte) (string, error) {
	args := m.Called(chunkHash, content)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorage) GetFileMetadata(fileHash string) (*model.FileMetadata, error) {
	args := m.Called(fileHash)
	return args.Get(0).(*model.FileMetadata), args.Error(1)
}

func (m *MockFileStorage) GetChunkData(chunkHash string) ([]byte, error) {
	args := m.Called(chunkHash)
	return args.Get(0).([]byte), args.Error(1)
}

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStorage) GetUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

type MockTokenHandler struct {
	mock.Mock
}

func (m *MockTokenHandler) CreateTokenPair(userID uint, username string) (*auth.TokenPair, error) {
	args := m.Called(userID, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenPair), args.Error(1)
}

func (m *MockTokenHandler) VerifyToken(token string) (*auth.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *MockTokenHandler) RefreshAccessToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}
