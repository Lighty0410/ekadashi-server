package grpc

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
)

type Server struct{}

type Service struct {
	controller *controller.Controller
}

func CreateServer(c *controller.Controller) *Service {
	return &Service{
		controller: c,
	}
}

func (s *Service) HandleRegistration(ctx context.Context, u *User) (*Response, error) {
	err := validateRequest(u)
	if err != nil {
		return nil, fmt.Errorf("cannot validate request: %v", err)
	}
	err = s.controller.RegisterUser(controller.User{Username: u.User, Password: u.Password})
	if err != nil {
		return nil, fmt.Errorf("gRPC request failed: %v", err)
	}
	return &Response{
		Response: "register succeed!",
	}, nil
}

func (s *Service) HandleLogin(ctx context.Context, u *User) (*Response, error) {
	err := validateRequest(u)
	if err != nil {
		return nil, fmt.Errorf("cannot validate request: %v", err)
	}
	session, err := s.controller.LoginUser(controller.User{Username: u.User, Password: u.Password})
	if err != nil {
		return nil, err
	}
	return &Response{
			Response: "register succeed!",
			Session:  &Session{Token: session.Token}},
		nil
}

func (s *Service) ShowEkadashi(ctx context.Context, u *User) (*Response, error) {
}

var validPassword = regexp.MustCompile(`^[a-zA-Z0-9=]+$`).MatchString
var validUsername = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func validateRequest(u *User) error {
	const minSymbols = 6
	if !validUsername(u.User) {
		return fmt.Errorf("field username contain latin characters and numbers without space only")
	}
	if !validPassword(u.Password) {
		return fmt.Errorf("field password contain latin characters and numbers without space only")
	}
	if len(u.User) < minSymbols {
		return fmt.Errorf("field username could not be less than 6 characters")
	}
	if len(u.Password) < minSymbols {
		return fmt.Errorf("field password could not be less than 6 characters")
	}
	return nil
}
