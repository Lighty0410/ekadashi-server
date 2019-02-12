package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/mongodb/mongo-go-driver/mongo"
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
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		return nil, fmt.Errorf("could not dial mongo: %v", err)
	}
	db := client.Database("ekadashi")
	return &Service{
		db: db,
	}, nil
}

func (s *Service) CreateIndex() error {
	var opt options.IndexOptions
	opt.SetExpireAfterSeconds(30)
	model := mongo.IndexModel{
		Options: &opt,
	}
	c := s.db.Collection("session")
	_, err := c.Indexes().CreateOne(context.Background(), model)
	if err != nil {
		return fmt.Errorf("cannont create an index")
	}
	fmt.Println("did this")
	return nil
}
