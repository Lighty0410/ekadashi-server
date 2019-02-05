package server

import (
	"encoding/json"
	"net/http"
)

func jsonError(w http.ResponseWriter, status int, err error) {
	jsonResponse(w, status, map[string]string{"reason": err.Error()})
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
