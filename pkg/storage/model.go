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

type SessionService interface {
	GetSession(hash string) (*Session, error)
	UpdateSession(session *Session) error
	CreateSession(u *Session) error
}

type EkadashiService interface {
	AddEkadashi(day *EkadashiDate) error
	NextEkadashi(day time.Time) (*EkadashiDate, error)
}

type UserService interface {
	AddUser(*User) error
	ReadUser(username string) (User, error)
	GetUsers() ([]string, error)
}
