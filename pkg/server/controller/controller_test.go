package controller

import (
	"fmt"
	"testing"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/crypto"
	"github.com/Lighty0410/ekadashi-server/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestController_RegisterUser(t *testing.T) {
	mockService := &storage.ServiceMock{
		AddUserFunc: func(in1 *storage.User) error {
			if in1.Name == "username" {
				return fmt.Errorf("duplicate key error collection")
			}
			if in1.Name == "somethingExceptable" {
				return fmt.Errorf("cannot add user")
			}
			return nil
		},
	}
	c := NewController(mockService)
	tt := []struct {
		casename      string
		user          User
		expectedError error
	}{
		{
			casename: "duplicate user",
			user: User{
				Username: "username",
				Password: "thisismypassword",
			},
			expectedError: ErrAlreadyExists,
		},
		{
			casename: "cannot add user",
			user: User{
				Username: "somethingExceptable",
				Password: "thisisordinady",
			},
			expectedError: fmt.Errorf("cannot add user to the database: cannot add user"),
		},
		{
			casename: "return nil",
			user: User{
				Username: "bestuserofalltime",
				Password: "qwerty123456",
			},
			expectedError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.casename, func(t *testing.T) {
			err := c.RegisterUser(tc.user)
			require.Equal(t, tc.expectedError, err)
		})
	}

}

func TestController_LoginUser(t *testing.T) {
	sessionCounter := 0
	mockService := &storage.ServiceMock{
		GetUserFunc: func(username string) (user storage.User, e error) {
			if username == "idontexist" {
				return user, storage.ErrUserNotFound
			}
			if username == "thisismongoDBnotme" {
				return user, fmt.Errorf("db crashed")
			}
			if username == "youauser" {
				hash, _ := crypto.GenerateHash("thisisnormpassword")
				user.PasswordHash = hash
				return user, nil
			}
			if username == "finallyitsfine" {
				hash, _ := crypto.GenerateHash("finepassword")
				user.PasswordHash = hash
				return user, nil
			}
			return user, nil
		},
		AddSessionFunc: func(in1 *storage.Session) error {
			if sessionCounter == 0 {
				sessionCounter++
				return fmt.Errorf("mongo error")
			}
			return nil
		},
	}
	tt := []struct {
		caseName      string
		user          User
		expectedError error
	}{
		{
			caseName: "user not found",
			user: User{
				Username: "idontexist",
				Password: "hehehepassword",
			},
			expectedError: ErrNotFound,
		},
		{
			caseName: "an error in mongoDB",
			user: User{
				Username: "thisismongoDBnotme",
				Password: "exceptionalpassword",
			},
			expectedError: fmt.Errorf("an error occurred in mongoDB during read user: db crashed"),
		},
		{
			caseName: "cannot create session",
			user: User{
				Username: "youauser",
				Password: "thisisnormpassword",
			},
			expectedError: fmt.Errorf("cannot create a session: mongo error"),
		},
		{
			caseName: "nil on return",
			user: User{
				Username: "finallyitsfine",
				Password: "finepassword",
			},
			expectedError: nil,
		},
	}
	c := NewController(mockService)
	for _, tc := range tt {
		t.Run(tc.caseName, func(t *testing.T) {
			_, err := c.LoginUser(tc.user)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestController_ShowEkadashi(t *testing.T) {
	nextEkadashi := 0
	mockService := &storage.ServiceMock{
		GetSessionFunc: func(hash string) (session *storage.Session, e error) {
			if hash == "zerotoken" {
				return nil, storage.ErrNoSession
			}
			if hash == "badauth" {
				return nil, fmt.Errorf("an error in mongodb")
			}
			session = new(storage.Session)
			session.Token = hash
			return session, nil
		},
		UpdateSessionFunc: func(in1 *storage.Session) error {
			if in1.Token == "cannotupdate" {
				return fmt.Errorf("bad session")
			}
			return nil
		},
		NextEkadashiFunc: func(in1 time.Time) (ekadashi *storage.Ekadashi, e error) {
			if nextEkadashi == 0 {
				nextEkadashi++
				return nil, fmt.Errorf("no ekadashi in dates")
			}
			ekadashi = new(storage.Ekadashi)
			ekadashi.Date = time.Date(2019, time.April, 16, 20, 30, 0, 0, time.UTC)
			return ekadashi, nil
		},
	}
	tt := []struct {
		caseName      string
		token         string
		expectedTime  time.Time
		expectedError error
	}{
		{
			caseName:      "user not found",
			token:         "zerotoken",
			expectedTime:  time.Time{},
			expectedError: fmt.Errorf("user not found"),
		},
		{
			caseName:      "cannot update user session",
			token:         "cannotupdate",
			expectedTime:  time.Time{},
			expectedError: fmt.Errorf("cannot check authentication: cannot update user session: bad session"),
		},
		{
			caseName:      "cannot check authentication",
			token:         "badauth",
			expectedTime:  time.Time{},
			expectedError: fmt.Errorf("cannot check authentication: an error in mongodb"),
		},
		{
			caseName:      "cannot get next ekadashi day",
			token:         "theresnoekadashi",
			expectedTime:  time.Time{},
			expectedError: fmt.Errorf("cannot get next ekadashi day: no ekadashi in dates"),
		},
		{
			caseName:      "return nil",
			token:         "ohthisisnicetoken",
			expectedTime:  time.Date(2019, time.April, 16, 20, 30, 0, 0, time.UTC),
			expectedError: nil,
		},
	}
	c := NewController(mockService)
	for _, tc := range tt {
		t.Run(tc.caseName, func(t *testing.T) {
			date, err := c.ShowEkadashi(tc.token)
			require.Equal(t, err, tc.expectedError)
			require.Equal(t, date, tc.expectedTime)
		})
	}
}
