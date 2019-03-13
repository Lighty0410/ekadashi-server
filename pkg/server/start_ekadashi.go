package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/server/moonapi"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

func (s *EkadashiServer) startEkadashi(ctx context.Context) error {
	ekadashi, _ := s.db.NextEkadashi(time.Now())
	if ekadashi.Date.IsZero() {
		days, err := moonapi.FillEkadashi()
		if err != nil {
			return fmt.Errorf("cannot fill ekadashi date: %v", err)
		}
		err = s.saveEkadashi(days)
		if err != nil {
			return err
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
				days, err := moonapi.FillEkadashi()
				if err != nil {
					log.Println("cannot fill ekadashi date: ", err)
					timer.Reset(time.Hour)
					continue
				}
				err = s.saveEkadashi(days)
				if err != nil {
					log.Println(err)
				}
				currentEkadashi, _ := s.db.NextEkadashi(time.Now())
				timer.Reset(currentEkadashi.Date.Sub(time.Now().Add(time.Hour * 24)))
			}
		}
	}()
	return nil
}

func (s *EkadashiServer) saveEkadashi(ekadashiDate []moonapi.SunMoon) error {
	for _, ekadashiDay := range ekadashiDate {
		err := s.db.AddEkadashi(&mongo.EkadashiDate{Date: ekadashiDay.Sun.RiseISO})
		if err != nil {
			return fmt.Errorf("cannot add date to the database: %v", err)
		}
	}
	return nil
}
