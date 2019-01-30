package mongo

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

// Service is used to interact with ekadashi storage.
type Service struct {
	db *mgo.Database
}

// NewService attempts to connect to MongoDB at localhost and
// if connection succeeds it returns Service ready to use.
func NewService() (*Service, error) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		return nil, fmt.Errorf("could not dial mongo: %v", err)
	}
	db := session.DB("ekadashi")
	return &Service{
		db: db,
	}, nil
}
