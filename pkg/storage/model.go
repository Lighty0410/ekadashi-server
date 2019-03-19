package storage

import (
	"fmt"
	"time"
)

// ErrUserNotFound is an error that returns if user is not found.
var ErrUserNotFound = fmt.Errorf("mongo: no documents in result")

// ErrNoSession is returned when session is not found.
var ErrNoSession = fmt.Errorf("session not found")

// Session contains an information about user session.
type Session struct {
	Token            string    `bson:"hash"`
	LastModifiedDate time.Time `bson:"modified"`
}

// ErrNoEkadashi is returned if there's no ekadashi dates in mongo.
var ErrNoEkadashi = fmt.Errorf("cannot find next ekadashi date")

// User contains an information about single user in a system.
type User struct {
	Name         string `bson:"name"`
	PasswordHash string `bson:"hash"`
}

// EkadashiDate is a structure that contains information about ekadashi date.
type Ekadashi struct {
	Date time.Time `bson:"date"`
}

// Service in the interface that is used to handle user's CRUD operations, user's session and ekadashi date.
type Service interface {
	GetSession(hash string) (*Session, error)
	UpdateSession(*Session) error
	AddSession(*Session) error
	AddEkadashi(*Ekadashi) error
	NextEkadashi(time.Time) (*Ekadashi, error)
	AddUser(*User) error
	GetUser(username string) (User, error)
}
