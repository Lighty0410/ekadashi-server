package server

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func generateHash(password string) (string, error) {
	cryptoHash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return "", fmt.Errorf("cannot generate hash: %v", err)
	}
	return string(cryptoHash), nil
}
func compareHash(hash string, password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), password)
}
