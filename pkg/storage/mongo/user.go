package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/storage"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// AddUser adds passed user into users collection.
func (s *Service) AddUser(u *storage.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c := s.db.Collection("users")
	_, err := c.InsertOne(ctx, u)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}

// GetUser retrieves an information from the database and compares it with a request.
func (s *Service) GetUser(username string) (storage.User, error) {
	var hash storage.User
	filter := bson.D{{Key: "name", Value: username}}
	err := s.db.Collection("users").FindOne(context.Background(), filter).Decode(&hash)
	if err == mongo.ErrNoDocuments {
		return hash, storage.ErrUserNotFound
	}
	if err != nil {
		return hash, fmt.Errorf("cannot search user, %v", err)
	}
	return hash, nil
}
