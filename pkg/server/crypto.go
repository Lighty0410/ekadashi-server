package server

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func generateHash(password string) (string, error) {
	cryptoHash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Println("Incorrect password") // Are there any requirements to return string/error in this case ?
	}
	return string(cryptoHash), err
}

func compareHash(hash string, password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), password)
}
