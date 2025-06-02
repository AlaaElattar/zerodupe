package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func HashPassword(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}
