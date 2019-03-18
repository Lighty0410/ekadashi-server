package main

import (
	"log"
	"os"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	internalHTTP "github.com/Lighty0410/ekadashi-server/pkg/server/http"
)

const ekadashiURL = "EKADASHI_MONGO_URL"

func main() {
	connectionURL := os.Getenv(ekadashiURL)
	if connectionURL == "" {
		log.Fatalf("Innapropriate %v variable for mongoDB connection", ekadashiURL)
	}
	mongoService, err := mongo.NewService(connectionURL)
	if err != nil {
		log.Fatalf("Could not create mongo service: %v", err)
	}
	newController := controller.NewController(mongoService)
	err = internalHTTP.NewServer(newController)
	if err != nil {
		log.Fatalf("Could not create ekadashi server: %v", err)
	}
}
