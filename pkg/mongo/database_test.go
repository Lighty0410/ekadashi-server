package mongo

import "testing"

func TestUserInsert(t *testing.T) {
	testService, err := NewService()
	if err != nil {
		t.Error(
			"Can't create a new service",
		)
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

	for _, table := range tt {
		t.Run(table.name, func(t *testing.T) {
			result := testService.AddUser(&table.user)
			if result != nil {
				t.Error(
					"For ", table.user,
					"Expected ", table.expectError,
					"Got ", result,
				)
			}
		})
	}
}
