package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// User contains an information about single user in a system.
type User struct {
	Name string `bson:"name"`
	Hash string `bson:"hash"`
}

// ErrUserNotFound is an error that returns if user is not found.
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

// CreateSession gets an information about session and insert it to database.
func (s *Service) CreateSession(u *Session) error {
	c := s.db.Collection("session")
	_, err := c.InsertOne(context.Background(), u)
	if err != nil {
		return fmt.Errorf("cannot create a session: %v", err)
	}
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
		return hash, fmt.Errorf("cannot search user, %v", err)
	}
	return hash, nil
}

// GetUsers gets an information about username of users.
func (s *Service) GetUsers() ([]string, error) {
	c := s.db.Collection("users")
	findOption := options.Find()
	findOption.SetProjection(bson.M{"name": 1})
	filter := bson.M{}
	cur, err := c.Find(context.Background(), filter, findOption)
	if err != nil {
		return nil, fmt.Errorf("cannot search users: %v", err)
	}
	var userList []string
	for cur.Next(context.Background()) {
		var u User
		err := cur.Decode(&u)
		if err != nil {
			return nil, fmt.Errorf("cannot decode user: %v", err)
		}
		fmt.Println(u.Name)
		userList = append(userList, u.Name)
	}
	return userList, err
}
