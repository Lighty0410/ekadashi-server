package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Session contains an information about user session.
type Session struct {
	Name             string    `bson:"name"`
	SessionHash      string    `bson:"hash"`
	LastModifiedDate time.Time `bson:"modified"`
}

var ErrNoSession = fmt.Errorf("mongo: User may exist but is unauthorized")

// GetSession receive information about user's hash and if succeed, returns Session structure.
func (s *Service) GetSession(hash string) (*Session, error) {
	var session Session
	c := s.db.Collection("session")
	filter := bson.D{{Key: "hash", Value: hash}}
	err := c.FindOne(context.Background(), filter).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNoSession
	}
	if err != nil {
		return nil, fmt.Errorf("could not find session, hash for this user does not exist: %v", err)
	}
	return &session, nil
}

// UpdateSession updates TTL index of current session.
func (s *Service) UpdateSession(session *Session) error {
	c := s.db.Collection("session")
	_, err := c.UpdateOne(context.Background(), bson.D{{Key: "hash", Value: session.SessionHash}}, bson.D{{
		"$set", bson.D{{
			Key: "modified", Value: session.LastModifiedDate,
		}},
	},
	})
	if err != nil {
		return fmt.Errorf("cannot update session: %v", err)
	}
	return nil
}
