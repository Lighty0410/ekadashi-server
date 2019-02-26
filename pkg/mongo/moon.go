package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// EkadashiDate is a structure that contains information about ekadashi date.
type EkadashiDate struct {
	Year  int        `bson:"year"`
	Month time.Month `bson:"month"`
	Day   int        `bson:"day"`
}

// AddEkadashi add ekadashi date to the database.
func (s *Service) AddEkadashi(day *EkadashiDate) error {
	c := s.db.Collection("ekadashi")
	_, err := c.InsertOne(context.Background(), day)
	if err != nil {
		return fmt.Errorf("cannot insert date to mongo DB: %v", err)
	}
	return nil
}

// SendEkadashi retrieves information from database and send it to another function.
func (s *Service) SendEkadashi() ([]EkadashiDate, error) {
	c := s.db.Collection("ekadashi")
	cur, err := c.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("cannot search any date: %v", err)
	}
	var ekadashiDay []EkadashiDate
	for cur.Next(context.Background()) {
		var ekadashiIteration EkadashiDate
		err := cur.Decode(&ekadashiIteration)
		if err != nil {
			return nil, fmt.Errorf("cannot decode date: %v", err)
		}
		ekadashiDay = append(ekadashiDay, ekadashiIteration)
	}
	return ekadashiDay, nil
}
