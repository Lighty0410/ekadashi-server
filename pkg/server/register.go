package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode"

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
		jsonError(w, http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}
	hashedPassword, err := generateHash(req.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("incorrect password: %v", err))
		return
	}
	err = s.db.AddUser(&mongo.User{
		Name: req.Username,
		Hash: hashedPassword,
	})
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("could not add user: %v", err))
		return
	}
	jsonResponse(w, http.StatusOK, nil)
}

func (req *loginRequest) validateRequest() error {
	IsLetter := regexp.MustCompile(`^[a-zA-Z1-9]+$`).MatchString
	var count int
	if !IsLetter(req.Username) || !IsLetter(req.Password) {
		return fmt.Errorf("field username should contains latin characters only")
	}
	for amount, symbol := range req.Username {
		if unicode.IsSpace(symbol) {
			return fmt.Errorf("field username cannot contains empty space")
		}
		count = amount
	}
	if count < 6 {
		return fmt.Errorf("field username cannot contains less than 6 character")
	}
	for amount, symbol := range req.Password {
		if unicode.IsSpace(symbol) {
			return fmt.Errorf("field password cannot contains empty space")
		}
		count = amount
	}
	if count < 6 {
		return fmt.Errorf("field username or password cannot contains less than 6 character")
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
		jsonError(w, http.StatusBadRequest, fmt.Errorf("%v", err))
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
	err = compareHash(user.Hash, []byte(req.Password))
	if err != nil {
		jsonError(w, http.StatusUnauthorized, fmt.Errorf("incorrect username or password: %v", err))
		return
	}
	userSession := &mongo.Session{
		Name:             req.Username,
		SessionHash:      generateToken(),
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
