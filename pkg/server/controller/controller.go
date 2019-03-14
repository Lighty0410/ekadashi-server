package controller

// TODO i don't like ekadashiHTTP name though. Gonna refactor it ASAP (when someone gonna helps me with it LOL)

import (
	"fmt"
	"strings"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/crypto"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

// This constants are defined to handle errors.
const (
	StatusConflict            = "StatusConflict"
	StatusInternalServerError = "StatusInternalServerError"
	StatusOK                  = "StatusOK"
	StatusBadRequest          = "StatusBadRequest"
	StatusUnauthorized        = "StatusUnauthorized"
)

// User contains information about a single user.
type User struct {
	Username string
	Password string
}

// Session contains information about user's session.
type Session struct {
	Name             string
	SessionHash      string
	LastModifiedDate time.Time
}

// RegisterUser is a method that register user and have an access to the mongoDB.
func (s *Controller) RegisterUser(u User) (string, error) {
	hashedPassword, err := crypto.GenerateHash(u.Password)
	if err != nil {
		return StatusInternalServerError, err
	}
	err = s.db.AddUser(&mongo.User{
		Name: u.Username,
		Hash: hashedPassword,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return StatusConflict, err
		}
		return StatusInternalServerError, err
	}
	return StatusOK, err
}

// LoginUser compares user's hash and password in the database.
// If succeed it add user's session to the database and returns it.
func (s *Controller) LoginUser(u User) (string, *Session, error) {
	user, err := s.db.ReadUser(u.Username)
	if err == mongo.ErrUserNotFound {
		return StatusUnauthorized, nil, fmt.Errorf("incorrect username or password: %v", err)
	}
	if err != nil {
		return StatusInternalServerError, &Session{}, fmt.Errorf("an error occurred in mongoDB: %v", err)
	}
	err = crypto.CompareHash(user.Hash, []byte(u.Password))
	if err != nil {
		return StatusUnauthorized, nil, fmt.Errorf("incorrect username or password: %v", err)
	}
	userSession := &mongo.Session{
		Name:             u.Username,
		SessionHash:      crypto.GenerateToken(),
		LastModifiedDate: time.Now(),
	}
	err = s.db.CreateSession(userSession)
	if err != nil {
		return StatusInternalServerError, nil, fmt.Errorf("cannot create a session: %v", err)
	}
	return StatusOK,
		&Session{Name: userSession.Name,
			SessionHash:      userSession.SessionHash,
			LastModifiedDate: userSession.LastModifiedDate}, // TODO don't like it. TOOOOOO LONG
		nil
}

// ShowEkadashi checks an existing session wi
func (s *Controller) ShowEkadashi(session string) (string, time.Time, error) { // TODO add commentaries. String or struct ? ReturnEkadashi, EkadashiDate, ShowEkadashi ?
	err := s.checkAuth(session)
	if err == mongo.ErrNoSession {
		return StatusUnauthorized, time.Now(), err
	}
	if err != nil {
		return StatusInternalServerError, time.Now(), fmt.Errorf("cannot check authentification: %v", err)
	}
	ekadashiDate, err := s.db.NextEkadashi(time.Now())
	if err != nil {
		return StatusInternalServerError, time.Now(), fmt.Errorf("cannot get next ekadashi day: %v", err)
	}
	return StatusOK, ekadashiDate.Date, nil // TODO time. time or struct ? hm
}

// checkAuth check current user's session.
// Return nil if succeed.
func (s *Controller) checkAuth(token string) error { // TODO where to place it ?
	session, err := s.db.GetSession(token)
	if err != nil {
		return err
	}
	session.LastModifiedDate = time.Now()
	err = s.db.UpdateSession(session)
	if err != nil {
		return err
	}
	return nil
}
