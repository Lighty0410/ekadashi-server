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
	SessionHash      string    `bson:"hash"`
	LastModifiedDate time.Time `bson:"modified"`
}

// User contains an information about single user in a system.
type User struct {
	Name string `bson:"name"`
	Hash string `bson:"hash"`
}

// EkadashiDate is a structure that contains information about ekadashi date.
type EkadashiDate struct {
	Date time.Time `bson:"date"`
}

// Service in the interface that uses to handle user's CRUD operations, user's session and ekadashi date.
type Service interface {
	GetSession(hash string) (*Session, error)
	UpdateSession(session *Session) error
	AddSession(u *Session) error
	AddEkadashi(day *EkadashiDate) error
	NextEkadashi(day time.Time) (*EkadashiDate, error)
	AddUser(*User) error
	GetUser(username string) (User, error)
}
