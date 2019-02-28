package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

type sunMoonResponse struct {
	Success bool      `json:"success"`
	Err     error     `json:"error"`
	Resp    []sunMoon `json:"response"`
}

type sunMoon struct {
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

func (s *EkadashiServer) fillEkadashi() error {
	accessID := os.Getenv(clientID)
	secretKey := os.Getenv(clientSecret)
	if accessID == "" || secretKey == "" {
		return fmt.Errorf("invalid accessID or secretkey value")
	}
	url := fmt.Sprintf("http://api.aerisapi.com/sunmoon/minsk,belarusmn?from=now&to=1month&limit=31&client_id=%s&client_secret=%s",
		accessID, secretKey)
	var moonPhase sunMoonResponse
	err := getJSON(url, &moonPhase)
	if err != nil {
		return fmt.Errorf("cannot get API server: %v", err)
	}
	if !moonPhase.Success {
		return fmt.Errorf("cannot succeed with API response: %v", moonPhase.Err)
	}

	filteredDate := ekadashiFilter(moonPhase.Resp)
	dayFilter := s.shiftEkadashi(filteredDate)
	err = s.saveEkadashi(dayFilter)
	if err != nil {
		return fmt.Errorf("an error occurred in mongoDB :%v", err)
	}
	return nil
}

func ekadashiFilter(sm []sunMoon) []sunMoon {
	const (
		newMoon  = "new moon"
		fullMoon = "full moon"
	)
	var ekadashiDate []sunMoon
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

func (s *EkadashiServer) shiftEkadashi(ekadashiDate []sunMoon) []sunMoon {
	var ekadashiDay []sunMoon
	for _, ekadashiFilter := range ekadashiDate {
		if ekadashiFilter.Moon.RiseISO.After(ekadashiFilter.Sun.RiseISO) {
			ekadashiFilter.Sun.RiseISO = ekadashiFilter.Sun.RiseISO.Add(time.Hour * 24)
		}
		ekadashiDay = append(ekadashiDay, ekadashiFilter)
	}
	return ekadashiDay
}

func (s *EkadashiServer) saveEkadashi(ekadashiDate []sunMoon) error {
	for _, ekadashiDay := range ekadashiDate {
		err := s.db.AddEkadashi(&mongo.EkadashiDate{Date: ekadashiDay.Sun.RiseISO})
		if err != nil {
			return fmt.Errorf("cannot add date to the database: %v", err)
		}
	}
	return nil
}
