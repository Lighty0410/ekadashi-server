package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// Service is used to interact with ekadashi storage.
type Service struct {
	db *mongo.Database
}

// NewService attempts to connect to MongoDB at localhost and
// if connection succeeds it returns Service ready to use.
func NewService(connectionURL string) (*Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, connectionURL)
	if err != nil {
		return nil, fmt.Errorf("could not dial mongo: %v", err)
	}
	db := client.Database("ekadashi")
	s := &Service{
		db: db,
	}
	err = s.createIndex()
	if err != nil {
		return s, fmt.Errorf("cannot create an index: %v", err)
	}
	return s, nil
}

// CreateIndex creates an index for collections.
func (s *Service) createIndex() error {
	var modifiedOpt, hashOpt options.IndexOptions
	modifiedOpt.SetExpireAfterSeconds(60 * 5)
	hashOpt.SetUnique(true)
	modifiedKey := bson.M{"modified": 1}
	hashKey := bson.M{"hash": 1}
	model := []mongo.IndexModel{
		{Keys: modifiedKey, Options: &modifiedOpt},
		{Keys: hashKey, Options: &hashOpt},
	}
	c := s.db.Collection("session")
	_, err := c.Indexes().CreateMany(context.Background(), model)
	if err != nil {
		return fmt.Errorf("cannont create session index: %v", err)
	}
	var usernameOpt options.IndexOptions
	usernameKey := bson.M{"username": 1}
	usernameOpt.SetUnique(true)
	userModel := mongo.IndexModel{
		Keys:    usernameKey,
		Options: &usernameOpt,
	}
	c = s.db.Collection("users")
	_, err = c.Indexes().CreateOne(context.Background(), userModel)
	if err != nil {
		return fmt.Errorf("cannot create index for 'users' collection: %v", err)
	}
	return nil
}
