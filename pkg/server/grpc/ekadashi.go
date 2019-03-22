package grpc

import (
	"context"
	"fmt"

	api "github.com/Lighty0410/ekadashi-server/pkg/server/grpc/api"
)

// ShowEkadashi send an endpoint request to the controller.
// If succeed returns ekadashi date.
func (s *Service) ShowEkadashi(ctx context.Context, u *api.Session) (*api.ShowEkadashiResponse, error) {
	if u.Token == "" {
		return nil, fmt.Errorf("cannot handle an empty token")
	}
	date, err := s.controller.ShowEkadashi(u.Token)
	if err != nil {
		return nil, err
	}
	ekadashiDate := date.Unix()
	return &api.ShowEkadashiResponse{Ekadashi: ekadashiDate}, nil
}
