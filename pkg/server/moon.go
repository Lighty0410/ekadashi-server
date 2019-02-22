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

func getJSON(url string, target interface{}) error {
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
	url := fmt.Sprintf("http://api.aerisapi.com/sunmoon/minsk,mn?from=now&to=1month&limit=31&client_id=%s&client_secret=%s",
		accessID, secretKey)
	moonPhase := sunMoonResponse{}
	err := getJSON(url, &moonPhase)
	if err != nil {
		err = fmt.Errorf("cannot get API server: %v", err)
		return err
	} else if !moonPhase.Success || moonPhase.Err != nil {
		err = fmt.Errorf("cannot get API server: %v", err)
		return err
	}
	filteredDate := ekadashiFilter(moonPhase.Resp)
	err = s.moonDateTransmit(filteredDate)
	if err != nil {
		err = fmt.Errorf("cannot insert ekadashi date todatabase: %v", err)
		return err
	}
	return nil
}

func ekadashiFilter(m []sunMoon) []sunMoon {
	const (
		newMoon  = "new moon"
		fullMoon = "full moon"
	)
	var ekadashiDate []sunMoon
	isNewMoon := false
	ekadashiDays := 0
	for _, date := range m {
		if date.Moon.Phase.Name == newMoon || date.Moon.Phase.Name == fullMoon {
			isNewMoon = true
		}
		if isNewMoon && !date.Moon.RiseISO.IsZero() {
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

func (s *EkadashiServer) moonDateTransmit(ekadashiDate []sunMoon) error {
	for _, ekadashiDay := range ekadashiDate {
		ekadashi := ekadashiDay.Sun.RiseISO
		if ekadashiDay.Moon.RiseISO.After(ekadashiDay.Sun.RiseISO) {
			ekadashi = ekadashi.Add(24 * time.Hour)
		}
		year, month, day := ekadashi.Date()
		err := s.db.AddMoonPhases(&mongo.EkadashiDate{Year: year, Month: month, Date: day})
		if err != nil {
			return fmt.Errorf("cannot add time to the database: %v", err)
		}
	}
	return nil
}
