package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/crypto"
	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

// ErrAlreadyExists is returned when such username already exists in the system.
var ErrAlreadyExists = fmt.Errorf("user already exists")

// ErrNotFound is returned when username doesn't exist or is incorrect.
var ErrNotFound = fmt.Errorf("user not found")

// User contains information about a single user.
type User struct {
	Username string
	Password string
}

// Session contains information about user's session.
type Session struct {
	Name             string
	Token            string
	LastModifiedDate time.Time
}

// Controller is an object that provides an access for the controller's functionality.
type Controller struct {
	db *mongo.Service
}

// CreateController creates a new instance for the controller.
func NewController(db *mongo.Service) *Controller {
	c := &Controller{
		db: db,
	}
	return c
}

// RegisterUser adds user and hashed password to the database.
// If succeed returns nil.
func (c *Controller) RegisterUser(u User) error {
	hashedPassword, err := crypto.GenerateHash(u.Password)
	if err != nil {
		return fmt.Errorf("cannot generate hash: %v", err)
	}
	err = c.db.AddUser(&mongo.User{
		Name: u.Username,
		Hash: hashedPassword,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return ErrAlreadyExists
		}
		return fmt.Errorf("cannot add user to the database: %v", err)
	}
	return nil
}

// LoginUser compares user's hash and password in the database.
// If succeed it add user's session to the database and returns it.
func (c *Controller) LoginUser(u User) (*Session, error) {
	user, err := c.db.ReadUser(u.Username)
	if err == mongo.ErrUserNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("an error occurred in mongoDB during read user: %v", err)
	}
	err = crypto.CompareHash(user.Hash, []byte(u.Password))
	if err != nil {
		return nil, ErrNotFound
	}
	userSession := &mongo.Session{
		Name:             u.Username,
		SessionHash:      crypto.GenerateToken(),
		LastModifiedDate: time.Now(),
	}
	err = c.db.CreateSession(userSession)
	if err != nil {
		return nil, fmt.Errorf("cannot create a session: %v", err)
	}
	return &Session{
		Name:             userSession.Name,
		Token:            userSession.SessionHash,
		LastModifiedDate: userSession.LastModifiedDate,
	}, nil
}

// ShowEkadashi checks an existing session.
// If succeed returns ekadashi date.
func (c *Controller) ShowEkadashi(sessionToken string) (time.Time, error) { //
	err := c.checkAuth(sessionToken)
	if err == mongo.ErrNoSession {
		return time.Time{}, ErrNotFound
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot check authentification: %v", err)
	}
	ekadashiDate, err := c.db.NextEkadashi(time.Now())
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot get next ekadashi day: %v", err)
	}
	return ekadashiDate.Date, nil
}

// checkAuth check current user's session.
// Return nil if succeed.
func (c *Controller) checkAuth(token string) error {
	session, err := c.db.GetSession(token)
	if err != nil {
		return fmt.Errorf("cannot get user session: %v", err)
	}
	session.LastModifiedDate = time.Now()
	err = c.db.UpdateSession(session)
	if err != nil {
		return fmt.Errorf("cannot update user session: %v", err)
	}
	return nil
}
