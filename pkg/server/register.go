package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// handleRegistration registers user in the system.
func (s *EkadashiServer) handleRegistration(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	var password string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("Can not decode the request: %v", err)) //"The specified password is not in the correct format. The password must be a non-empty string." What do you think about this ?
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
	var req registerRequest
	var hash string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("Can not decode the request.:  %v", err)) // Which kind of problem is it ? Is it the server-side problem or client one ?
		return
	}
	hash, err = s.db.ReadUser(req.Username)
	if err != nil {
		jsonError(w, http.StatusForbidden, fmt.Errorf("A user with the specified username and password combination does not exist in the system.: %v", err))
		return
	}
	err = compareHash(hash, []byte(req.Password))
	if err != nil {
		jsonError(w, http.StatusForbidden, fmt.Errorf("A user with the specified username and password combination does not exist in the system.: %v", err))
		return
	}
	jsonResponse(w, http.StatusOK, nil)
}
