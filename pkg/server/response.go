package server

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, payload interface{}) {

	json.NewEncoder(w).Encode(payload)
}
