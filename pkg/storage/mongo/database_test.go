package mongo

import (
	"os"
	"testing"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		user        storage.User
		expectError error
	}{
		{
			name:        "empty password",
			user:        storage.User{Name: "Greatestmateofalltime", PasswordHash: ""},
			expectError: nil,
		},
		{
			name:        "empty name",
			user:        storage.User{Name: "", PasswordHash: ""},
			expectError: nil,
		},
		{
			name:        "empty password",
			user:        storage.User{Name: "Leva", PasswordHash: ""},
			expectError: nil,
		},
		{
			name:        "ASCII symobols as a string",
			user:        storage.User{Name: "@!#@!#", PasswordHash: "123213"},
			expectError: nil,
		},
		{
			name:        "casual database info",
			user:        storage.User{Name: "Mesropyan", PasswordHash: "SecretKey"},
			expectError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := testService.AddUser(&tc.user)
			require.NoError(t, err)
			user, err := testService.GetUser(tc.user.Name)
			require.NoError(t, err)
			assert.Equal(t, user, tc.user, tc.name)
		})
	}
}

func TestService_NextEkadashiAndAddEkadashi(t *testing.T) {
	connectionURL := os.Getenv("EKADASHI_MONGO_URL")
	require.NotEmpty(t, connectionURL)
	testService, err := NewService(connectionURL)
	require.NoError(t, err)
	tt := []struct {
		name         string
		userDate     time.Time
		date         []storage.Ekadashi
		expectedDate time.Time
	}{
		{
			name: "dateSince before the first date",
			userDate: time.Date(
				2009, 11, 17, 20, 34, 58, 0, time.UTC),
			date: []storage.Ekadashi{
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
			date: []storage.Ekadashi{
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
			date: []storage.Ekadashi{
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
				require.NoError(t, testService.AddEkadashi(&date))
			}
			ekadashiDate, err := testService.NextEkadashi(tc.userDate)
			if assert.NoError(t, err) {
				assert.Equal(t, ekadashiDate.Date.UTC(), tc.expectedDate.UTC())
			}
		})
	}
}
