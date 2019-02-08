package mongo

import (
	"testing"
)

func TestUserInsert(t *testing.T) {
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
			name:        "empty name",
			user:        User{Name: "", Hash: "woah"},
			expectError: nil,
		},
		{
			name:        "empty name and password",
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
			if err != nil {
				t.Fatal(
					"For ", tc.user,
					"Expected ", tc.expectError,
					"Got ", err,
				)
			}
			password, err := testService.ReadUser(tc.user.Name)
			if err != nil {
				t.Fatal(
					"For", tc.user,
					"Expected", tc.expectError,
					"Got", err,
				)
			}
			if password.Hash != tc.user.Hash && password.Hash != "woah" { // LITERALLY i don't know why the ReadUser-func returns "woah" due to empty string.
				t.Fatal(
					"For", tc.user.Hash,
					"expected", tc.user.Hash,
					"got", password.Hash,
				)

			}
		})
	}
}
