package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

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
func generateToken() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
