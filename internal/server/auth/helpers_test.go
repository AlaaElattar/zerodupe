package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashAndSaltPasswordAndVerifyPassword(t *testing.T) {
	password := "mySecretPassword"

	hashed, err := HashAndSaltPassword([]byte(password))
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.Greater(t, len(hashed), 5)

	assert.True(t, VerifyPassword(hashed, password))

	assert.False(t, VerifyPassword(hashed, "wrongPassword"))

	hashed2, err := HashAndSaltPassword([]byte(password))
	assert.NoError(t, err)
	assert.NotEqual(t, hashed, hashed2)
	assert.True(t, VerifyPassword(hashed2, password))
}
