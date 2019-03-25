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
func (s *Service) Register(ctx context.Context, u *api.RegisterRequest) (*api.Empty, error) {
	err := validateRequest(u.User)
	if err != nil {
		return nil, fmt.Errorf("cannot validate the request: %v", err)
	}
	err = s.controller.RegisterUser(controller.User{Username: u.User.Name, Password: u.User.Password})
	if err != nil {
		return nil, fmt.Errorf("could not register user: %v", err)
	}
	return &api.Empty{}, nil
}

// HandleRegistration validate login request and sends it to the controller.
// If succeed it login user in the system and returns successful response.
func (s *Service) Login(ctx context.Context, u *api.LoginRequest) (*api.LoginResponse, error) {
	err := validateRequest(u.User)
	if err != nil {
		return nil, fmt.Errorf("cannot validate the request: %v", err)
	}
	session, err := s.controller.LoginUser(controller.User{Username: u.User.Name, Password: u.User.Password})
	if err != nil {
		return nil, fmt.Errorf("could not login user: %v", err)
	}
	return &api.LoginResponse{
		Response: &api.Session{Token: session.Token}}, nil
}
