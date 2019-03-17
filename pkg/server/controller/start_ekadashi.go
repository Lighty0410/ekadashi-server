package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/ekadashi"
	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

func (c *Controller) StartEkadashi(ctx context.Context) error {
	actual, err := c.db.NextEkadashi(time.Now())
	if actual.Date.IsZero() || err == mongo.ErrNoEkadashi {
		ekadashiDate, err := ekadashi.FillEkadashi()
		if err != nil {
			return fmt.Errorf("cannot fill ekadashi date: %v", err)
		}
		err = c.saveEkadashi(ekadashiDate)
		if err != nil {
			return fmt.Errorf("cannot save ekadashi: %v", err)
		}
	}
	actual, err = c.db.NextEkadashi(time.Now())
	if err != nil {
		return fmt.Errorf("cannot get next ekadashi day: %v", err)
	}
	go func() {
		timer := time.NewTimer(actual.Date.Sub(time.Now().Add(time.Hour * 24)))
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				ekadashDate, err := ekadashi.FillEkadashi()
				if err != nil {
					log.Println("cannot fill ekadashi date: ", err)
					timer.Reset(time.Hour)
					continue
				}
				err = c.saveEkadashi(ekadashDate)
				if err != nil {
					log.Println("cannot save ekadashi: %v", err)
				}
				currentEkadashi, err := c.db.NextEkadashi(time.Now())
				if err != nil {
					log.Println("cannot get next ekadashi day: %v", err)
				}
				timer.Reset(currentEkadashi.Date.Sub(time.Now().Add(time.Hour * 24)))
			}
		}
	}()
	return nil
}

func (c *Controller) saveEkadashi(ekadashiDate []ekadashi.Date) error {
	for _, ekadashiDay := range ekadashiDate {
		err := c.db.AddEkadashi(&mongo.EkadashiDate{Date: ekadashiDay.Sun.RiseISO})
		if err != nil {
			return fmt.Errorf("cannot add date to the database: %v", err)
		}
	}
	return nil
}
