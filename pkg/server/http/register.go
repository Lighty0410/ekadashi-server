package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
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
	err = s.controller.RegisterUser(controller.User{Username: req.Username, Password: req.Password})
	if err != nil {
		switch err {
		case controller.ErrAlreadyExists:
			jsonError(w, http.StatusConflict, err)
			return
		default:
			jsonError(w, http.StatusInternalServerError, err)
			return
		}
	}
	jsonResponse(w, http.StatusOK, nil)
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
	session, err := s.controller.LoginUser(controller.User{Username: req.Username, Password: req.Password})
	if err != nil {
		switch err {
		case controller.ErrNotFound:
			jsonError(w, http.StatusUnauthorized, err)
			return
		default:
			jsonError(w, http.StatusInternalServerError, err)
			return
		}
	}
	cookie := http.Cookie{
		Name:  "sesssion_token",
		Value: session.Token,
	}
	http.SetCookie(w, &cookie)
	jsonResponse(w, http.StatusOK, nil)
}

var validPassword = regexp.MustCompile(`^[a-zA-Z0-9=]+$`).MatchString
var validUsername = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

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
