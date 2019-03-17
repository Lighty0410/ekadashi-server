package mongo

import (
	"context"
	"fmt"

	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

var ErrNoEkadashi = fmt.Errorf("cannot find next ekadashi date")

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

// NextEkadashi retrieves information about the last ekadashi date from the database.
func (s *Service) NextEkadashi(day time.Time) (*EkadashiDate, error) {
	var ekadashiDay EkadashiDate
	searchOpt := options.FindOneOptions{}
	searchOpt.Sort = bson.M{"date": 1}
	c := s.db.Collection("ekadashi")
	err := c.FindOne(context.Background(), bson.D{{
		Key: "date", Value: bson.D{{
			Key: "$gt", Value: day.Add(-24 * time.Hour),
		}},
	}}, &searchOpt).Decode(&ekadashiDay)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, nil
		}
		return &ekadashiDay, ErrNoEkadashi
	}
	return &ekadashiDay, nil
}
