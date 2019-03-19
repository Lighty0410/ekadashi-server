package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/storage"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// AddUser adds passed user into users collection.
func (s *Service) Add(u *storage.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c := s.db.Collection("users")
	_, err := c.InsertOne(ctx, u)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}

// ReadUser retrieves an information from the database and compares it with a request.
func (s *Service) Read(username string) (storage.User, error) {
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

// GetUsers gets an information about username of users.
func (s *Service) Get() ([]string, error) {
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
		var u storage.User
		err := cur.Decode(&u)
		if err != nil {
			return nil, fmt.Errorf("cannot decode user: %v", err)
		}
		fmt.Println(u.Name)
		userList = append(userList, u.Name)
	}
	return userList, err
}
