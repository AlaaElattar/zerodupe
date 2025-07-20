package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashAndSaltPassword hashes the given password using bcrypt
func HashAndSaltPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// VerifyPassword checks if the given password matches the hashed one using bcrypt
func VerifyPassword(hashedPassword []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
}
