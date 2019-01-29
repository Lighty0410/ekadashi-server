package server

import (
	"encoding/json"
	"net/http"
)

//Error returns an error if request's not valid
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, map[string]string{"error": message})
}

//JSON marshalling string to json
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
