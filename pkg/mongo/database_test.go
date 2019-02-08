package mongo

import (
	"testing"
)

func TestAddAndReadUser(t *testing.T) {
	testService, err := NewService()
	if err != nil {
		t.Error(
			"Can't create a new service",
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
					"For: ", tc.expectError,
					"Expected: ", tc.expectError,
					"Got: ", err,
				)
			}
			if user != tc.user {
				t.Fatal(
					"For: ", tc.expectError,
					"Expected: ", tc.expectError,
					"Got: ", err,
				)
			}
			if user != tc.user {
				t.Fatal(
					"For: ", tc.user,
					"Expected: ", tc.user,
					"Got: ", user,
				)
			}
		})
	}
}
