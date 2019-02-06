package mongo

import (
	"context"
	"fmt"
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

func (s *Service) ReadUser(username string) (string, error) {
	var result User
	filter := bson.D{{"name", username}}
	err := s.db.Collection("users").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return fmt.Sprintf("	A user with the specified username and password combination does not exist in the system."), err //Is this necessary ?
	}
	return result.Hash, err
}
