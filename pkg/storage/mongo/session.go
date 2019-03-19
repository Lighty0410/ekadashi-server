package mongo

import (
	"context"
	"fmt"

	"github.com/Lighty0410/ekadashi-server/pkg/storage"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// CreateSession gets an information about session and insert it to database.
func (s *Service) CreateSession(u *storage.Session) error {
	c := s.db.Collection("session")
	_, err := c.InsertOne(context.Background(), u)
	if err != nil {
		return fmt.Errorf("cannot create a session: %v", err)
	}
	return nil
}

// GetSession receive information about user's hash and if succeed, returns Session structure.
func (s *Service) GetSession(hash string) (*storage.Session, error) {
	var session storage.Session
	c := s.db.Collection("session")
	filter := bson.D{{Key: "hash", Value: hash}}
	err := c.FindOne(context.Background(), filter).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, storage.ErrNoSession
	}
	if err != nil {
		return nil, fmt.Errorf("could not find session: %v", err)
	}
	return &session, nil
}

// UpdateSession updates TTL index of current session.
func (s *Service) UpdateSession(session *storage.Session) error {
	c := s.db.Collection("session")
	_, err := c.UpdateOne(context.Background(), bson.D{{Key: "hash", Value: session.SessionHash}}, bson.D{{
		Key: "$set", Value: bson.D{{
			Key: "modified", Value: session.LastModifiedDate,
		}},
	},
	})
	if err != nil {
		return fmt.Errorf("cannot update session: %v", err)
	}
	return nil
}
