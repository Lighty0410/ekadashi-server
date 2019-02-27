package mongo

import (
	"context"
	"fmt"
	"time"
)

// EkadashiDate is a structure that contains information about ekadashi date.
type EkadashiDate struct {
	Date time.Time `bson:"date"`
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
