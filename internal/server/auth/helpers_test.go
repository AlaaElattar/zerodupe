package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGenerateSalt(t *testing.T) {
	t.Run("Test GenerateSalt generates a random salt", func(t *testing.T) {
		salt1, err := GenerateSalt()
		assert.NoError(t, err)
		salt2, err := GenerateSalt()
		assert.NoError(t, err)
		assert.NotEqual(t, salt1, salt2)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("Test HashPassword generates a consistent hash for the same password and salt", func(t *testing.T) {
		password := "test"
		salt, err := GenerateSalt()
		assert.NoError(t, err)
		hash1 := HashPassword(password, salt)
		hash2 := HashPassword(password, salt)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("Test HashPassword generates a different hash for a different password", func(t *testing.T) {
		password1 := "test"
		password2 := "test1"
		salt, err := GenerateSalt()
		assert.NoError(t, err)
		hash1 := HashPassword(password1, salt)
		hash2 := HashPassword(password2, salt)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Test HashPassword generates a different hash for a different salt", func(t *testing.T) {
		password := "test"
		salt1, err := GenerateSalt()
		assert.NoError(t, err)
		salt2, err := GenerateSalt()
		assert.NoError(t, err)
		hash1 := HashPassword(password, salt1)
		hash2 := HashPassword(password, salt2)
		assert.NotEqual(t, hash1, hash2)
	})
}
