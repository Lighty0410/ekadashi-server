package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	"github.com/Lighty0410/ekadashi-server/pkg/server/grpc"
	"github.com/Lighty0410/ekadashi-server/pkg/server/http"
	"github.com/Lighty0410/ekadashi-server/pkg/storage/mongo"
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
	httpServer, err := http.NewServer(":9000", newController)
	if err != nil {
		log.Fatalf("Could not create ekadashi server: %v", err)
	}
	grpcServer, err := grpc.NewGrpcServer(":50051", newController)
	if err != nil {
		log.Fatalf("Could not create gRPC server: %v", err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	log.Printf("Shutting down due to signal: %v", sig)
	err = httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("cannot shutdown the server: %v", err)
	}
	grpcServer.GracefulStop()
}
