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
		for {
			timer := time.NewTimer(time.Minute * 5)
			<-timer.C
			if time.Now().Day() > ekadashi.Date.Day() {
				err := s.db.UpdateEkadashi(time.Now().AddDate(0, 0, -1))
				if err != nil {
					log.Println("cannot update ekadashi: ", err)
				}
				err = s.fillEkadashi()
				if err != nil {
					log.Println("cannot fill ekadashi date: ", err)
				}
			}
		}
	}()
	return nil
}
