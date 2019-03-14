package controller

import "github.com/Lighty0410/ekadashi-server/pkg/mongo"

// Controller is an object that provides an access for the controller's functionality.
type Controller struct {
	db *mongo.Service
}

// CreateController creates a new instance for the controller.
func NewController(db *mongo.Service) *Controller {
	c := &Controller{
		db: db,
	}
	return c
}
