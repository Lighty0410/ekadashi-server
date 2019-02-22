package mongo

import (
	"context"
	"fmt"
	"time"
)

// EkadashiDate is a structure that contains information about ekadashi date.
type EkadashiDate struct {
	Year  int        `bson:"year"`
	Month time.Month `bson:"month"`
	Date  int        `bson:"day"`
}

// AddMoonPhases add ekadashi date to the database.
func (s *Service) AddMoonPhases(days *EkadashiDate) error {
	c := s.db.Collection("ekadashi")
	_, err := c.InsertOne(context.Background(), days)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
