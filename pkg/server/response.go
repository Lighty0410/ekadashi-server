package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonError(w http.ResponseWriter, status int, err error) {
	jsonResponse(w, status, map[string]string{"reason": err.Error()})
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		fmt.Errorf("Cannot encode the interface: %v", payload)
		return
	}
}
