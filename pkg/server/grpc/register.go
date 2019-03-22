package grpc

import (
	"context"
	"fmt"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	api "github.com/Lighty0410/ekadashi-server/pkg/server/grpc/api"
)

// Service is a struct that contains controller field.
// Every gRPC method is defined by this structure.
type Service struct {
	controller *controller.Controller
}

// NewService creates a new instance for service.
func NewService(c *controller.Controller) *Service {
	return &Service{
		controller: c,
	}
}

// HandleRegistration validate register request and sends it to the controller.
// If succeed it register user in the system and returns successful response.
func (s *Service) HandleRegistration(ctx context.Context, u *api.User) (*api.Response, error) {
	err := validateRequest(u)
	if err != nil {
		return nil, fmt.Errorf("cannot validate request: %v", err)
	}
	err = s.controller.RegisterUser(controller.User{Username: u.User, Password: u.Password})
	if err != nil {
		return nil, fmt.Errorf("gRPC request failed: %v", err)
	}
	return &api.Response{
		Response: "register succeed!",
	}, nil
}

// HandleRegistration validate login request and sends it to the controller.
// If succeed it register user in the system and returns successful response.
func (s *Service) HandleLogin(ctx context.Context, u *api.User) (*api.Response, error) {
	err := validateRequest(u)
	if err != nil {
		return nil, fmt.Errorf("cannot validate request: %v", err)
	}
	session, err := s.controller.LoginUser(controller.User{Username: u.User, Password: u.Password})
	if err != nil {
		return nil, err
	}
	return &api.Response{
			Response: "login succeed!",
			Session:  &api.Session{Name: "session_token", Token: session.Token}},
		nil
}
