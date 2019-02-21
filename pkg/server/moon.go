package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

type MoonPhases struct {
	Success bool                `json:"success"`
	Err     error               `json:"error"`
	Resp    []MoonPhaseResponse `json:"response"`
}

type MoonPhaseResponse struct {
	Moon Moon `json:"moon"`
	Sun  Sun  `json:"sun"`
}

type Sun struct {
	RiseISO time.Time `json:"riseISO"`
}

type Moon struct {
	Rise    int       `json:"rise"`
	RiseISO time.Time `json:"riseISO"`
	Phase   Phase     `json:"phase"`
}

type Phase struct {
	Name string `json:"name"`
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
func (s *EkadashiServer) moonAPI(w http.ResponseWriter, _ *http.Request) {
	url := "http://api.aerisapi.com/sunmoon/minsk,mn?from=now&to=1month&limit=31&client_id=6foK0Hivo3udMxkoWGpL5&client_secret=xrgfvYtd3x9MY0XGK4Q7NvUoIRY5W0U2WEhTOJDC"
	moonPhase := &MoonPhases{}
	err := getJSON(url, &moonPhase)
	if err != nil {
		jsonError(w, http.StatusBadGateway, fmt.Errorf("cannot get API server: %v", err))
		return
	}
	err = s.ekadashiFilter(moonPhase.Resp)
}

func (s *EkadashiServer) ekadashiFilter(m []MoonPhaseResponse) error {
	var ekadashiDate []MoonPhaseResponse
	isNewMoon := false
	ekadashiDays := 0
	for _, date := range m {

		if date.Moon.Phase.Name == "new moon" || date.Moon.Phase.Name == "full moon" {
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
	err := s.moonDateTransmit(ekadashiDate)
	if err != nil {
		return fmt.Errorf("cannot transmit date to database: %v", err)
	}
	return nil
}

func (s *EkadashiServer) moonDateTransmit(ekadashiDate []MoonPhaseResponse) error {
	for _, ekadashiDay := range ekadashiDate {
		if ekadashiDay.Moon.RiseISO.Before(ekadashiDay.Sun.RiseISO) {
			year, month, day := ekadashiDay.Sun.RiseISO.Date()
			err := s.db.AddMoonPhases(&mongo.EkadashiDate{Year: year, Month: month, Day: day})
			if err != nil {
				return fmt.Errorf("cannot add time to the database: %v", err)
			}
		} else {
			ekadashi := ekadashiDay.Sun.RiseISO.Add(24 * time.Hour)
			year, month, day := ekadashi.Date()
			err := s.db.AddMoonPhases(&mongo.EkadashiDate{Year: year, Month: month, Day: day})
			if err != nil {
				return fmt.Errorf("cannot add time to the database: %v", err)
			}
		}
	}
	return nil
}
