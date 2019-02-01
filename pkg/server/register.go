package server

import (
	"encoding/json"
	"fmt"
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
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("could not decode request: %v", err))
		return
	}
	err = s.db.AddUser(&mongo.User{
		Name:     req.Username,
		Password: req.Password,
	})
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("could not add user: %v", err))
		return
	}
	jsonResponse(w, http.StatusOK, nil)
}

// here's handlers
