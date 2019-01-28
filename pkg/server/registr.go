package server

import (
	"encoding/json"
	"net/http"
)

// AuthName is a struct of users that's gonna connect to the server
type AuthName struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Server is a struct wich we can to handle to
type Server struct {
	http.Server
	shutDown chan bool
	reqCount uint32
}

var users []AuthName

// Login authorize users.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var user AuthName
	json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}
