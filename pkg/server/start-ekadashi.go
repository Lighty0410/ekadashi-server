package server

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (s *EkadashiServer) startEkadashi(ctx context.Context) error {
	ekadashi, _ := s.db.NextEkadashi(time.Now())
	if ekadashi.Date.IsZero() {
		err := s.fillEkadashi()
		if err != nil {
			return fmt.Errorf("cannot fill ekadashi date: %v", err)
		}
	}
	go func() {
		timer := time.NewTimer(ekadashi.Date.Sub(time.Now().Add(time.Hour * 24)))
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				err := s.fillEkadashi()
				if err != nil {
					log.Println("cannot fill ekadashi date: ", err)
					timer.Reset(time.Hour)
					continue
				}
				currentEkadashi, _ := s.db.NextEkadashi(time.Now())
				timer.Reset(currentEkadashi.Date.Sub(time.Now().Add(time.Hour * 24)))
			}
		}
	}()
	return nil
}
