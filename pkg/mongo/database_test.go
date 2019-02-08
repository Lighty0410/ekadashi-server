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
			if err != tc.expectError {
				t.Fatal(
					"For: ", tc.user,
					"Expected: ", tc.expectError,
					"Got: ", err,
				)
			}
		})
	}
}

func TestReadUser(t *testing.T) {
	testService, err := NewService()
	if err != nil {
		t.Error("cannot create a new service")
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
			user:        User{Name: "Leva", Hash: "SecretKey"},
			expectError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := testService.AddUser(&tc.user)
			if err != tc.expectError {
				t.Fatal(
					"For: ", tc.user,
					"Expected: ", tc.expectError,
					"Got: ", err,
				)
			}
			user, err := testService.ReadUser(tc.user.Name)
			if err != tc.expectError {
				t.Fatal(
					"For: ", tc.user,
					"Expected ", tc.expectError,
					"Got: ", err,
				)
			}
			if user.Name != tc.user.Name {
				t.Fatal(
					"For: ", tc.user,
					"Expected: ", tc.user,
					"Got: ", user.Name,
				)
			}
		})
	}
}
