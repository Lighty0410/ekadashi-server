package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/crypto"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	loginRequest
}

// handleRegistration registers user in the system.
func (s *EkadashiServer) handleRegistration(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("can not decode the request: %v", err))
		return
	}
	err = req.validateRequest()
	if err != nil {
		jsonError(w, http.StatusBadRequest, err)
		return
	}
	hashedPassword, err := crypto.GenerateHash(req.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("incorrect password: %v", err))
		return
	}
	err = s.db.AddUser(&mongo.User{
		Name: req.Username,
		Hash: hashedPassword,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			jsonError(w, http.StatusConflict, fmt.Errorf("user already exists"))
			return
		}
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("could not add user: %v", err))
		return
	}
	jsonResponse(w, http.StatusOK, nil)
}

var validPassword = regexp.MustCompile(`^[a-zA-Z1-9=]+$`).MatchString
var validUsername = regexp.MustCompile(`^[a-zA-Z1-9]+$`).MatchString

func (req *loginRequest) validateRequest() error {
	const minSymbols = 6
	if !validUsername(req.Username) {
		return fmt.Errorf("field username contain latin characters and numbers without space only")
	}
	if !validPassword(req.Password) {
		return fmt.Errorf("field password contain latin characters and numbers without space only")
	}
	if len(req.Username) < minSymbols {
		return fmt.Errorf("field username could not be less than 6 characters")
	}
	if len(req.Password) < minSymbols {
		return fmt.Errorf("field password could not be less than 6 characters")
	}
	return nil
}

// checkAuth check current user's session.
// Return nil if succeed.
func (s *EkadashiServer) checkAuth(token string) error {
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

// handleLogin retrieves information about login request.
// If login succeed it assigns cookie to user.
func (s *EkadashiServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("can not decode the request: %v", err))
		return
	}
	err = req.validateRequest()
	if err != nil {
		jsonError(w, http.StatusBadRequest, err)
		return
	}
	user, err := s.db.ReadUser(req.Username)
	if err == mongo.ErrUserNotFound {
		jsonError(w, http.StatusUnauthorized, fmt.Errorf("incorrect username or password: %v", err))
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("an error occurred in mongoDB: %v", err))
		return
	}
	err = crypto.CompareHash(user.Hash, []byte(req.Password))
	if err != nil {
		jsonError(w, http.StatusUnauthorized, fmt.Errorf("incorrect username or password: %v", err))
		return
	}
	userSession := &mongo.Session{
		Name:             req.Username,
		SessionHash:      crypto.GenerateToken(),
		LastModifiedDate: time.Now(),
	}
	err = s.db.CreateSession(userSession)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("cannot create a session: %v", err))
	}
	cookie := http.Cookie{
		Name:  "session_token",
		Value: userSession.SessionHash,
	}
	http.SetCookie(w, &cookie)
	jsonResponse(w, http.StatusOK, nil)
}
