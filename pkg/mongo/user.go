package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// User contains information about a single user in a system.
type User struct {
	Name string `bson:"name"`
	Hash string `bson:"hash"`
}

type Session struct {
	SessionID string `bson:"sessionid"`
}

// ErrUserNotFound is an error that returns if user is not found
var ErrUserNotFound = fmt.Errorf("mongo: no documents in result")

// AddUser adds passed user into users collection.
func (s *Service) AddUser(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c := s.db.Collection("users")
	_, err := c.InsertOne(ctx, u)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}

func (s *Service) EnsureSession(session *Session) error {
	s.CreateIndex()
	c := s.db.Collection("session")
	_, err := c.InsertOne(context.Background(), session)
	if err != nil {
		return fmt.Errorf("cannot add session")
	}
	fmt.Println("OK")
	return nil
}

// ReadUser retrieves an information from the database and compares it with a request.
func (s *Service) ReadUser(username string) (User, error) {
	var hash User
	filter := bson.D{{Key: "name", Value: username}}
	err := s.db.Collection("users").FindOne(context.Background(), filter).Decode(&hash)
	if err == mongo.ErrNoDocuments {
		return hash, ErrUserNotFound
	}
	if err != nil {
		return hash, fmt.Errorf("cannot search user")
	}
	return hash, nil
}
