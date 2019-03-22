package grpc

import (
	"fmt"
	"regexp"

	api "github.com/Lighty0410/ekadashi-server/pkg/server/grpc/api"
)

var validPassword = regexp.MustCompile(`^[a-zA-Z0-9=]+$`).MatchString
var validUsername = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func validateRequest(u *api.User) error {
	const minSymbols = 6
	if !validUsername(u.Name) {
		return fmt.Errorf("field username contain latin characters and numbers without space only")
	}
	if !validPassword(u.Name) {
		return fmt.Errorf("field password contain latin characters and numbers without space only")
	}
	if len(u.Password) < minSymbols {
		return fmt.Errorf("field username could not be less than 6 characters")
	}
	if len(u.Password) < minSymbols {
		return fmt.Errorf("field password could not be less than 6 characters")
	}
	return nil
}
