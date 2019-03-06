package server

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func convertStringToTime(t *testing.T, str string) time.Time {
	dateTime, err := time.Parse("2006-01-02T15:04:05-07:00", str)
	if err != nil {
		t.Log(err)
		return time.Time{}
	}
	return dateTime
}

func TestEkadashi(t *testing.T) {
	testService := &EkadashiServer{}
	file, err := os.Open("../../test-data/ekadashi-data.json")
	if err != nil {
		t.Fatal("cannot open file: ", err)
	}
	defer file.Close()
	var testMoonDate []sunMoon
	err = json.NewDecoder(file).Decode(&testMoonDate)
	if err != nil {
		t.Fatal("cannot decode file: ", err)
	}
	date := ekadashiFilter(testMoonDate)
	date = testService.shiftEkadashi(date)
	expectedData := []sunMoon{
		{Sun: sun{RiseISO: convertStringToTime(t, "2019-01-17T07:45:05-06:00")}},
		{Sun: sun{RiseISO: convertStringToTime(t, "2019-01-31T07:32:06-06:00")}},
	}
	for i, tc := range expectedData {
		if !tc.Sun.RiseISO.Equal(date[i].Sun.RiseISO) {
			t.Fatalf("For: %v\nExpected: %v\nGot: %v",
				tc.Sun.RiseISO, tc.Sun.RiseISO, date[i].Sun.RiseISO)
		}
	}
}
