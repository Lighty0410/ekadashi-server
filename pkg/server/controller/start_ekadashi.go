package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/ekadashi"
	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

// FillEkadashi is a goroutine that autofill ekadashi dates via getEkadashi().
// If succeed it adds ekadashi date to the database.
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
				}
				timer.Reset(nextEkadashi.Date.Sub(time.Now().Add(time.Hour * 24)))
			}
		}
	}()
	return nil
}

func (c *Controller) getEkadashi() (*mongo.EkadashiDate, error) {
	ek, err := c.db.NextEkadashi(time.Now())
	if err != nil && err != mongo.ErrNoEkadashi {
		return nil, err
	}
	if err == mongo.ErrNoEkadashi {
		dates, err := ekadashi.NextMonth()
		if err != nil {
			return nil, fmt.Errorf("cannot fill ekadashi date: %v", err)
		}
		err = c.saveEkadashi(dates)
		if err != nil {
			return nil, fmt.Errorf("cannot save ekadashi: %v", err)
		}
		ek, err = c.db.NextEkadashi(time.Now())
		if err == mongo.ErrNoEkadashi {
			return nil, err
		}
	}
	return ek, nil
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
