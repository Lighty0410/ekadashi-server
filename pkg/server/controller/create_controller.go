package controller

import "github.com/Lighty0410/ekadashi-server/pkg/mongo"

type Controller struct {
	db *mongo.Service
}

func CreateController(db *mongo.Service) *Controller {
	c := &Controller{
		db: db,
	}
	return c
}
