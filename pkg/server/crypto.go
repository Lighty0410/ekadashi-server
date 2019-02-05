package server

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func generateHash(password string) (string, error) {
	cryptoHash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Println("Incorrect password")
	}
	return string(cryptoHash), err
}

func compareHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
}
