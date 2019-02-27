package server

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (s *EkadashiServer) startEkadashi(ctx context.Context) error {
	ekadashi, _ := s.db.RetrieveEkadashi(time.Now())
	if ekadashi.Date.IsZero() {
		err := s.fillEkadashi()
		if err != nil {
			return fmt.Errorf("cannot fill ekadashi date: %v", err)
		}
	}
	go func() {
		timer := time.NewTimer(time.Now().Add(time.Hour * 24).Sub(ekadashi.Date))
		for {
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				err := s.fillEkadashi()
				if err != nil {
					log.Println("cannot fill ekadashi date: ", err)
				}
				currentEkadashi, _ := s.db.RetrieveEkadashi(time.Now())
				timer.Reset(time.Now().Add(time.Hour * 24).Sub(currentEkadashi.Date))
			}
		}
	}()
	return nil
}
