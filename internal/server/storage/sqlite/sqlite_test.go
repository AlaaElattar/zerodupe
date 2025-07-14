package sqlite

// import (
// 	"testing"
// 	"zerodupe/internal/server/model"

// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/gorm"
// )

// func TestCreateUser(t *testing.T) {
// 	t.Run("Test CreateUser creates a new user", func(t *testing.T) {
// 		sqliteStorage := NewSqliteStorageMock()
// 		user := &model.User{
// 			Username: "test",
// 		}
// 		err := sqliteStorage.CreateUser(user, "test")

// 		assert.NoError(t, err)
// 		assert.Equal(t, "test", user.Username)
// 		assert.NotEmpty(t, user.Password)
// 		assert.NotEqual(t, user.Password, "test")
// 		assert.NotEmpty(t, user.Salt)

// 		retrievedUser, err := sqliteStorage.GetUserByUsername("test")
// 		assert.NoError(t, err)
// 		assert.Equal(t, user.Username, retrievedUser.Username)
// 		assert.Equal(t, user.Password, retrievedUser.Password)
// 		assert.Equal(t, user.Salt, retrievedUser.Salt)
// 	})

// 	t.Run("Test CreateUser returns error for existing user", func(t *testing.T) {
// 		sqliteStorage := NewSqliteStorageMock()

// 		user1 := &model.User{
// 			Username: "test",
// 		}
// 		err := sqliteStorage.CreateUser(user1, "test")
// 		assert.NoError(t, err)
// 		assert.Equal(t, "test", user1.Username)
// 		assert.NotEmpty(t, user1.Password)
// 		assert.NotEqual(t, user1.Password, "test")
// 		assert.NotEmpty(t, user1.Salt)

// 		user2 := &model.User{
// 			Username: "test",
// 		}
// 		err = sqliteStorage.CreateUser(user2, "test")
// 		assert.Error(t, err)
// 		assert.Equal(t, gorm.ErrDuplicatedKey, err)
// 	})
// }
// func TestGetUserByUsername(t *testing.T) {
// 	t.Run("Test GetUserByUsername returns user for existing user", func(t *testing.T) {
// 		sqliteStorage := NewSqliteStorageMock()
// 		user := &model.User{
// 			Username: "test",
// 		}
// 		err := sqliteStorage.CreateUser(user, "test")
// 		assert.NoError(t, err)

// 		retrievedUser, err := sqliteStorage.GetUserByUsername("test")
// 		assert.NoError(t, err)
// 		assert.Equal(t, user.Username, retrievedUser.Username)
// 		assert.Equal(t, user.Password, retrievedUser.Password)
// 		assert.Equal(t, user.Salt, retrievedUser.Salt)
// 	})

// 	t.Run("Test GetUserByUsername returns error for non-existing user", func(t *testing.T) {
// 		sqliteStorage := NewSqliteStorageMock()
// 		_, err := sqliteStorage.GetUserByUsername("test")
// 		assert.Error(t, err)
// 		assert.Equal(t, gorm.ErrRecordNotFound, err)
// 	})
// }
