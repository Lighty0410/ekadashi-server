package mongo

import (
	"context"
	"fmt"

	"github.com/Lighty0410/ekadashi-server/pkg/storage"

	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// ErrNoEkadashi is returned if there's no ekadashi dates in mongo.
var ErrNoEkadashi = fmt.Errorf("cannot find next ekadashi date")

// AddEkadashi add ekadashi date to the database.
func (s *Service) AddEkadashi(day *storage.EkadashiDate) error {
	c := s.db.Collection("ekadashi")
	_, err := c.InsertOne(context.Background(), day)
	if err != nil {
		return fmt.Errorf("cannot insert date to mongo DB: %v", err)
	}
	return nil
}

// NextEkadashi retrieves information about the last ekadashi date from the database.
func (s *Service) NextEkadashi(day time.Time) (*storage.EkadashiDate, error) {
	var ekadashiDay storage.EkadashiDate
	searchOpt := options.FindOneOptions{}
	searchOpt.Sort = bson.M{"date": 1}
	c := s.db.Collection("ekadashi")
	err := c.FindOne(context.Background(), bson.D{{
		Key: "date", Value: bson.D{{
			Key: "$gt", Value: day.Add(-24 * time.Hour),
		}},
	}}, &searchOpt).Decode(&ekadashiDay)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoEkadashi
		}
		return nil, err
	}
	return &ekadashiDay, nil
}
