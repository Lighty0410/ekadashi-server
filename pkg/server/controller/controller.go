package controller

// TODO i don't like ekadashiHTTP name though. Gonna refactor it ASAP (when someone gonna helps me with it LOL)

import (
	"fmt"
	"strings"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
	"github.com/Lighty0410/ekadashi-server/pkg/server/helper"
)

const (
	StatusConflict            = "StatusConflict"
	StatusInternalServerError = "StatusInternalServerError"
	StatusOK                  = "StatusOK"
	StatusBadRequest          = "StatusBadRequest"
	StatusUnauthorized        = "StatusUnauthorized"
)

type User struct { // TODO commentaries
	Username string
	Password string
}

type Session struct {
	Name             string
	SessionHash      string
	LastModifiedDate time.Time
}

type EkadashiDate struct {
	Date time.Time
}

func (s *Controller) RegisterUser(u User) (string, error) { // TODO commentaries
	hashedPassword, err := helper.GenerateHash(u.Password)
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

func (s *Controller) LoginUser(u User) (string, Session, error) { // TODO add commentaries.
	user, err := s.db.ReadUser(u.Username)
	if err == mongo.ErrUserNotFound {
		return StatusUnauthorized, Session{}, fmt.Errorf("incorrect username or password: %v", err)
	}
	if err != nil {
		return StatusInternalServerError, Session{}, fmt.Errorf("an error occurred in mongoDB: %v", err)
	}
	err = helper.CompareHash(user.Hash, []byte(u.Password))
	if err != nil {
		return StatusUnauthorized, Session{}, fmt.Errorf("incorrect username or password: %v", err)
	}
	userSession := &mongo.Session{
		Name:             u.Username,
		SessionHash:      helper.GenerateToken(),
		LastModifiedDate: time.Now(),
	}
	err = s.db.CreateSession(userSession)
	if err != nil {
		return StatusInternalServerError, Session{}, fmt.Errorf("cannot create a session: %v", err)
	}
	return StatusOK,
		Session{Name: userSession.Name,
			SessionHash:      userSession.SessionHash,
			LastModifiedDate: userSession.LastModifiedDate}, // TODO don't like it. TOOOOOO LONG
		nil
}

func (s *Controller) ShowEkadashi(session string) (string, EkadashiDate, error) { // TODO add commentaries. String or struct ?
	err := s.checkAuth(session)
	if err == mongo.ErrNoSession {
		return StatusUnauthorized, EkadashiDate{}, err
	}
	if err != nil {
		return StatusInternalServerError, EkadashiDate{}, fmt.Errorf("cannot check authentification: %v", err)
	}
	ekadashiDate, err := s.db.NextEkadashi(time.Now())
	if err != nil {
		return StatusInternalServerError, EkadashiDate{}, fmt.Errorf("cannot get next ekadashi day: %v", err)
	}
	return StatusOK, EkadashiDate{Date: ekadashiDate.Date}, nil // TODO time. time or struct ? hm
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
