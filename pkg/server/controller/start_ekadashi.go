package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/ekadashi"
	"github.com/Lighty0410/ekadashi-server/pkg/storage"
)

// FillEkadashi is a goroutine that autofills ekadashi dates.
// If succeed it adds ekadashi dates to the database.
func (c *Controller) FillEkadashi(ctx context.Context) error {
	actualEkadashi, err := c.getEkadashi()
	if err != nil {
		return err
	}
	go func() {
		timer := time.NewTimer(actualEkadashi.Date.Sub(time.Now().Add(time.Hour * 24)))
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				nextEkadashi, err := c.getEkadashi()
				if err != nil {
					log.Println(err)
					timer.Reset(time.Hour)
					continue
				}
				timer.Reset(nextEkadashi.Date.Sub(time.Now().Add(time.Hour * 24)))
			}
		}
	}()
	return nil
}

func (c *Controller) getEkadashi() (*storage.Ekadashi, error) {
	ek, err := c.service.NextEkadashi(time.Now())
	if err != nil && err != storage.ErrNoEkadashi {
		return nil, err
	}
	if err == storage.ErrNoEkadashi {
		dates, err := ekadashi.NextMonth()
		if err != nil {
			return nil, fmt.Errorf("cannot fill ekadashi date: %v", err)
		}
		err = c.saveEkadashi(dates)
		if err != nil {
			return nil, fmt.Errorf("cannot save ekadashi: %v", err)
		}
		ek, err = c.service.NextEkadashi(time.Now())
		if err == storage.ErrNoEkadashi {
			return nil, err
		}
	}
	return ek, nil
}

func (c *Controller) saveEkadashi(ekadashiDate []ekadashi.Date) error {
	for _, ekadashiDay := range ekadashiDate {
		err := c.service.AddEkadashi(&storage.Ekadashi{Date: ekadashiDay.Sun.RiseISO})
		if err != nil {
			return fmt.Errorf("cannot add date to the database: %v", err)
		}
	}
	return nil
}
