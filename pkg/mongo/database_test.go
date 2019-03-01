package mongo

import (
	"os"
	"testing"
	"time"
)

func TestAddAndReadUser(t *testing.T) {
	connectionURL := os.Getenv("EKADASHI_MONGO_URL")
	if connectionURL == "" {
		t.Error(
			"inappropriate environment variable",
		)
		return
	}
	testService, err := NewService(connectionURL)
	if err != nil {
		t.Error(
			"can't create a new service",
		)
		return
	}
	tt := []struct {
		name        string
		user        User
		expectError error
	}{
		{
			name:        "empty password",
			user:        User{Name: "Greatestmateofalltime", Hash: ""},
			expectError: nil,
		},
		{
			name:        "empty name",
			user:        User{Name: "", Hash: ""},
			expectError: nil,
		},
		{
			name:        "empty password",
			user:        User{Name: "Leva", Hash: ""},
			expectError: nil,
		},
		{
			name:        "ASCII symobols as a string",
			user:        User{Name: "@!#@!#", Hash: "123213"},
			expectError: nil,
		},
		{
			name:        "casual database info",
			user:        User{Name: "Mesropyan", Hash: "SecretKey"},
			expectError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := testService.AddUser(&tc.user)
			if err != tc.expectError {
				t.Fatal(
					"For: ", tc.user,
					"\nExpected: ", tc.expectError,
					"\nGot: ", err,
				)
			}
			user, err := testService.ReadUser(tc.user.Name)
			if err != tc.expectError {
				t.Fatal(
					"For: ", tc.expectError,
					"\nExpected: ", tc.expectError,
					"\nGot: ", err,
				)
			}
			if user != tc.user {
				t.Fatal(
					"For: ", tc.user,
					"\nExpected: ", tc.user,
					"\nGot: ", user,
				)
			}
		})
	}
}

func TestService_NextEkadashiAndAddEkadashi(t *testing.T) {
	connectionURL := os.Getenv("MONGO_EKADASHI_URL")
	if connectionURL == "" {
		t.Error("incorrect environment variable")
		return
	}
	testService, err := NewService(connectionURL)
	if err != nil {
		t.Fatalf("cannot connect to database: %v", err)

	}
	tt := []struct {
		name         string
		userDate     time.Time
		date         []EkadashiDate
		expectedDate time.Time
		expectErr    error
	}{
		{
			name: "dateSince before the first date",
			userDate: time.Date(
				2009, 11, 17, 20, 34, 58, 0, time.UTC),
			date: []EkadashiDate{
				{Date: time.Date(
					2009, 11, 23, 20, 34, 58, 0, time.UTC)},
				{Date: time.Date(
					2009, 11, 26, 20, 34, 58, 0, time.UTC)},
				{Date: time.Date(
					2009, 11, 30, 20, 34, 58, 0, time.UTC)}},
			expectedDate: time.Date(
				2009, 11, 23, 20, 34, 58, 0, time.UTC),
		},
		{
			name: "dateSince before the second date",
			userDate: time.Date(
				2009, 11, 25, 20, 34, 58, 0, time.UTC),
			date: []EkadashiDate{
				{Date: time.Date(
					2009, 11, 23, 20, 34, 58, 0, time.UTC)},
				{Date: time.Date(
					2009, 11, 26, 20, 34, 58, 0, time.UTC)},
				{Date: time.Date(
					2009, 11, 30, 20, 34, 58, 0, time.UTC)}},
			expectedDate: time.Date(
				2009, 11, 26, 20, 34, 58, 0, time.UTC),
		},
		{
			name: "day to day",
			userDate: time.Date(
				2009, 11, 30, 20, 34, 58, 0, time.UTC),
			date: []EkadashiDate{
				{Date: time.Date(
					2009, 11, 23, 20, 34, 58, 0, time.UTC)},
				{Date: time.Date(
					2009, 11, 26, 20, 34, 58, 0, time.UTC)},
				{Date: time.Date(
					2009, 11, 30, 20, 34, 58, 0, time.UTC)}},
			expectedDate: time.Date(
				2009, 11, 30, 20, 34, 58, 0, time.UTC),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, date := range tc.date {
				err := testService.AddEkadashi(&date)
				if err != nil {
					t.Fatalf("an error occurred in database: %v", err)
				}
			}
			ekadashiDate, err := testService.NextEkadashi(tc.userDate)
			if err != tc.expectErr {
				t.Error("an error occurred in database: ", err)
			}
			if !ekadashiDate.Date.UTC().Equal(tc.expectedDate.UTC()) {
				t.Fatal(
					"For: ", tc.name, " ", tc.userDate,
					"\nexpected: ", tc.expectedDate,
					"\n     got: ", ekadashiDate.Date)
			}
		})
	}
}
