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
	var test = []User{
		User{Name: "a", Password: "wow"},
		User{Name: "", Password: "woah"},
		User{Name: "", Password: ""},
		User{Name: "@!#@!#", Password: "123213"},
		User{Name: "e", Password: "ouch"},
	}
	for _, pair := range test {
		result := testService.AddUser(&pair)
		if result != nil {
			t.Error(
				"For ", pair.Name,
				"Expected ", pair.Password,
				"Got ", u.Password, u.Name,
			)
		}
	}
}
