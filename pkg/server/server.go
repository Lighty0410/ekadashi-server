package server

import (
	"encoding/json"
	"net/http"
	"testsmt/pkg/mongo"
)

//Registration registrate users and insert user's info in MongoDB
func Registration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user mongo.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}
	if err := mongo.CreateUser(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	respondWithJSON(w, user)
}
