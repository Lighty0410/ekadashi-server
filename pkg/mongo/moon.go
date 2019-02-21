package mongo

import (
	"context"
	"fmt"
	"time"
)

type EkadashiDate struct {
	Year  int        `bson:"year"`
	Month time.Month `bson:"month"`
	Day   int        `bson:"day"`
}

func (s *Service) AddMoonPhases(phases *EkadashiDate) error {
	c := s.db.Collection("phases")
	_, err := c.InsertOne(context.Background(), phases)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
