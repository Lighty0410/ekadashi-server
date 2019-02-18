package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MoonPhases struct {
	Success bool       `json:"success"`
	Err     error      `json:"error"`
	Resp    []Response `json:"response"`
}

type Response struct {
	Timestamp   int    `json:"timestamp"`
	DateTimeISO string `json:"dateTimeISO"`
	Code        int    `json:"code"`
	Name        string `json:"name"`
}

func getJSON(url string, target interface{}) error {
	myClient := &http.Client{}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(&target)
}
func moonAPI(w http.ResponseWriter, _ *http.Request) {
	url := "http://api.aerisapi.com/sunmoon/moonphases?limit=4&client_id=6foK0Hivo3udMxkoWGpL5&client_secret=xrgfvYtd3x9MY0XGK4Q7NvUoIRY5W0U2WEhTOJDC"
	moonPhase := &MoonPhases{}
	err := getJSON(url, &moonPhase)
	if err != nil {
		jsonError(w, http.StatusBadGateway, fmt.Errorf("cannot get API server: %v", err))
	}
	jsonResponse(w, http.StatusOK, moonPhase.Resp)
}
