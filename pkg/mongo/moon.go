package mongo

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
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

// SendEkadashi retrieves information from database and send it to another function.
func (s *Service) RetrieveEkadashi(day time.Time) (EkadashiDate, error) {
	c := s.db.Collection("ekadashi")
	cur, err := c.Find(context.Background(), bson.D{{
		Key: "date", Value: bson.D{{
			Key: "$gt", Value: day,
		}},
	}})
	if err != mongo.ErrNoDocuments {
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

func (s *Service) UpdateEkadashi(date time.Time) error {
	c := s.db.Collection("ekadashi")
	_, err := c.DeleteMany(context.Background(), bson.D{{
		Key: "date", Value: bson.D{{
			Key: "$lt", Value: date,
		}},
	}})
	if err != nil {
		return fmt.Errorf("cannot delete date in database :%v", err)
	}
	return nil
}
