package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	var password string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("can not decode the request: %v", err))
		return
	}
	password, err = generateHash(req.Password)
	if err != nil {
		log.Println("Incorrect password")
		return
	}
	fmt.Println(req.Username)
	err = s.db.AddUser(&mongo.User{
		Name: req.Username,
		Hash: password,
	})
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("could not add user: %v", err))
		return
	}
	jsonResponse(w, http.StatusOK, nil)
}

func (s *EkadashiServer) handleLogin(w http.ResponseWriter, r *http.Request) { //login
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("can not decode the request.:  %v", err))
		return
	}
	user, status, err := s.db.ReadUser(req.Username)
	if err != nil {
		jsonError(w, status, fmt.Errorf("incorrect username or password: %v", err))
		return
	} else if status != http.StatusUnauthorized && status != http.StatusOK && nil != err {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("an error occurred in mongoDB: %v", err))
	}
	err = compareHash(user.Hash, []byte(req.Password))
	if err != nil {
		jsonError(w, http.StatusForbidden, fmt.Errorf("incorrect username or password: %v", err))
		return
	}
	jsonResponse(w, http.StatusOK, nil)
}
