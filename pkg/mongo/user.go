package mongo

import "fmt"

// User contains information about a single user in a system.
type User struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

// AddUser adds passed user into users collection.
func (s *Service) AddUser(u *User) error {
	c := s.db.C("users")
	err := c.Insert(u)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}
