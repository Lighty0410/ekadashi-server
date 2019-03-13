package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	"github.com/Lighty0410/ekadashi-server/pkg/server/ekadashihttp"
)

type EkadashiServer struct { // TODO does this name suits properly in this project/cases idk ?
	db *mongo.Service
}

func NewEkadashiServer(db *mongo.Service) error {
	s := EkadashiServer{
		db: db,
	}
	c := controller.CreateController(db)
	err := initHTTP(c)
	if err != nil {
		return err
	}
	err = s.initEkadashi()
	if err != nil {
		return err
	}
	return nil
}

func initHTTP(c *controller.Controller) error { // TODO init ? Rly ? And again i don't like it. What's the better name ?
	ekadashiServer, err := ekadashihttp.NewHttpServer(c)
	if err != nil {
		log.Fatalf("Could not create ekadashi server: %v", err)
	}
	server := &http.Server{
		Addr:    ":9000",
		Handler: ekadashiServer,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Could not listen: %v", err)
		}
	}()
	sig := <-stop
	log.Printf("Shutting down due to signal: %v", sig)
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Println("Cannot shutdown the server")
	}
	return nil
}

func (s *EkadashiServer) initEkadashi() error { // TODO again fucking init. It gonna pisses me off cuz of my own mindset. The second one: we already have startEkadashi :<
	err := s.startEkadashi(context.Background())
	if err != nil {
		return fmt.Errorf("cannot fill ekadashiAPI: %v", err)
	}
	return nil
}
