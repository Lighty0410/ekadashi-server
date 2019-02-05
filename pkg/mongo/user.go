package mongo

import (
	"context"
	"fmt"
	"time"
)

// User contains information about a single user in a system.
type User struct {
	Name string `bson:"name"`
	Hash string `bson:"hash"`
}

// AddUser adds passed user into users collection.
func (s *Service) AddUser(u *User) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	c := s.db.Collection("users")
	_, err := c.InsertOne(ctx, u) // err := c.Insert(u)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}

/*
func (s *Service) ReadUser() error {
	u := &User{}
	err := s.db.C("users").Find(u.Name)
	if err != nil {
		log.Println("incorrect collection")
	}
	return nil
}
*/
