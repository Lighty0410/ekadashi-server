package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
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

// LastEkadashi retrieves information about the last ekadashi date from the database.
func (s *Service) LastEkadashi(day time.Time) (EkadashiDate, error) {
	c := s.db.Collection("ekadashi")
	cur, err := c.Find(context.Background(), bson.D{{
		Key: "date", Value: bson.D{{
			Key: "$gt", Value: day,
		}},
	}})
	if err != nil {
		return EkadashiDate{}, fmt.Errorf("cannot find an existing file: %v", err)
	}
	var ekadashiDay EkadashiDate
	for cur.Next(context.Background()) {
		err := cur.Decode(&ekadashiDay)
		if err != nil {
			return EkadashiDate{}, fmt.Errorf("cannot decode date: %v", err)
		}
		break
	}
	return ekadashiDay, nil
}
