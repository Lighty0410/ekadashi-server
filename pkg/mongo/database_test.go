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
			user:        User{Name: "", Password: "woah"},
			expectError: nil,
		},
		{
			name:        "empty name and password",
			user:        User{Name: "", Password: ""},
			expectError: nil,
		},
		{
			name:        "empty password",
			user:        User{Name: "Leva", Password: ""},
			expectError: nil,
		},
		{
			name:        "ASCII symobols as a string",
			user:        User{Name: "@!#@!#", Password: "123213"},
			expectError: nil,
		},
		{
			name:        "casual database info",
			user:        User{Name: "Mesropyan", Password: "SecretKey"},
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
		})
	}
}
