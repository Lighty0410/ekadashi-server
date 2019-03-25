package grpc

import (
	"context"
	"fmt"

	api "github.com/Lighty0410/ekadashi-server/pkg/server/grpc/api"
)

// ShowEkadashi send an endpoint request to the controller.
// If succeed returns ekadashi date.
func (s *Service) ShowEkadashi(ctx context.Context, u *api.ShowEkadashiRequest) (*api.ShowEkadashiResponse, error) {
	if u.Request.Token == "" {
		return nil, fmt.Errorf("auth token is required")
	}
	date, err := s.controller.ShowEkadashi(u.Request.Token)
	if err != nil {
		return nil, fmt.Errorf("cannot get ekadashi date from gRPC: %v", err)
	}
	ekadashiDate := date.Unix()
	return &api.ShowEkadashiResponse{Ekadashi: ekadashiDate}, nil
}
