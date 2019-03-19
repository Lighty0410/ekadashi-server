package storage

import (
	"time"
)

// Session contains an information about user session.
type Session struct {
	Name             string    `bson:"name"`
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

type Sessioner interface {
	GetSession(hash string) (*Session, error)
	UpdateSession(session *Session) error
	CreateSession(u *Session) error
}

type Ekadasher interface {
	AddEkadashi(day *EkadashiDate) error
	NextEkadashi(day time.Time) (*EkadashiDate, error)
}

type Userer interface {
	AddUser(*User) error
	ReadUser(username string) (User, error)
	GetUsers() ([]string, error)
}
