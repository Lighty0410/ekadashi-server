package ekadashi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type sunMoonResponse struct {
	Success bool   `json:"success"`
	Err     error  `json:"error"`
	Resp    []Date `json:"response"`
}

type Date struct {
	Moon moon `json:"moon"`
	Sun  sun  `json:"sun"`
}

type sun struct {
	RiseISO time.Time `json:"riseISO"`
}

type moon struct {
	RiseISO time.Time `json:"riseISO"`
	Phase   phase     `json:"phase"`
}

type phase struct {
	Name string `json:"name"`
}

const (
	clientID     = "CLIENT_ID"
	clientSecret = "CLIENT_SECRET"
)

//nolint:errcheck
func getJSON(url string, target interface{}) (err error) {
	r, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot get url: %v", err)
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(&target)
}

func FillEkadashi() ([]Date, error) {
	accessID := os.Getenv(clientID)
	secretKey := os.Getenv(clientSecret)
	if accessID == "" || secretKey == "" {
		return nil, fmt.Errorf("invalid accessID or secretkey value")
	}
	url := fmt.Sprintf("http://api.aerisapi.com/sunmoon/minsk,belarusmn?from=now&to=1month&limit=31&client_id=%s&client_secret=%s",
		accessID, secretKey)
	var moonPhase sunMoonResponse
	err := getJSON(url, &moonPhase)
	if err != nil {
		return nil, fmt.Errorf("cannot get API server: %v", err)
	}
	if !moonPhase.Success {
		return nil, fmt.Errorf("cannot succeed with API response: %v", moonPhase.Err)
	}

	filteredDate := ekadashiFilter(moonPhase.Resp)
	dayFilter := shiftEkadashi(filteredDate)
	return dayFilter, nil
}

func ekadashiFilter(sm []Date) []Date {
	const (
		newMoon  = "new moon"
		fullMoon = "full moon"
	)
	var ekadashiDate []Date
	isNewMoon := false
	ekadashiDays := 0
	for _, date := range sm {
		if date.Moon.Phase.Name == newMoon || date.Moon.Phase.Name == fullMoon {
			isNewMoon = true
		}
		if isNewMoon {
			ekadashiDays++
		}
		if ekadashiDays == 11 {
			ekadashiDate = append(ekadashiDate, date)
			isNewMoon = false
			ekadashiDays = 0
		}
	}
	return ekadashiDate
}

func shiftEkadashi(ekadashiDate []Date) []Date {
	var ekadashiDay []Date
	for _, ekadashiFilter := range ekadashiDate {
		if ekadashiFilter.Moon.RiseISO.After(ekadashiFilter.Sun.RiseISO) {
			ekadashiFilter.Sun.RiseISO = ekadashiFilter.Sun.RiseISO.Add(time.Hour * 24)
		}
		ekadashiDay = append(ekadashiDay, ekadashiFilter)
	}
	return ekadashiDay
}
