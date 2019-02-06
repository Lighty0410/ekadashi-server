package mongo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// User contains information about a single user in a system.
type User struct {
	Name string `bson:"name"`
	Hash string `bson:"hash"`
}

// AddUser adds passed user into users collection.
func (s *Service) AddUser(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c := s.db.Collection("users")
	_, err := c.InsertOne(ctx, u) // err := c.Insert(u)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}

// ReadUser retrieves an information from the database and compares it with a request.
func (s *Service) ReadUser(username string) (User, int, error) {
	var result User
	filter := bson.D{{"name", username}}
	err := s.db.Collection("users").FindOne(context.Background(), filter).Decode(&result)
	fmt.Println(result.Name, result.Hash)
	switch err != nil {
	case err == errors.New("mongo: no documents in result"):
		return result, http.StatusUnauthorized, err
	case err == errors.New("Registry cannot be nil"):
		return result, http.StatusUnauthorized, err
	case err != nil:
		return result, http.StatusInternalServerError, err
	default:
		return result, http.StatusInternalServerError, err
	}
}
