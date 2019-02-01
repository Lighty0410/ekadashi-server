package mongo

import "testing"

func TestUserInsert(t *testing.T) {
	testService, err := NewService()
	u := &User{}
	if err != nil {
		t.Error(
			"Can't create a new service",
		)
	}
	tt := []struct {
		name        string
		user        User
		expectError string
	}{
		{
			name: "all ok",
			user: User{Name: "a", Password: "wow"},
		},
	}
	tt := []User{
		User{Name: "", Password: "woah"},
		User{Name: "", Password: ""},
		User{Name: "leva", Password: ""},
		User{Name: "@!#@!#", Password: "123213"},
		User{Name: "e", Password: "ouch"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := testService.AddUser(&tc.user)
			if result != nil {
				t.Error(
					"For ", pair.Name,
					"Expected ", pair.Password,
					"Got ", u.Password, u.Name,
				)
			}
		})
	}
}
