package grpc

import (
	"context"
	"fmt"

	api "github.com/Lighty0410/ekadashi-server/pkg/server/grpc/api"
)

// ShowEkadashi send an endpoint request to the controller.
// If succeed returns ekadashi date.
func (s *Service) ShowEkadashi(ctx context.Context, u *api.User) (*api.Response, error) {
	if u.Auth.Name != "session_token" {
		return nil, fmt.Errorf("incorrect session name")
	}
	date, err := s.controller.ShowEkadashi(u.Auth.Token)
	if err != nil {
		return nil, err
	}
	ekadashiDate := date.Unix()
	return &api.Response{Ekadashi: ekadashiDate}, nil
}
