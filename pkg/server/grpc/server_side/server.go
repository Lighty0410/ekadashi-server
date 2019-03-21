package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	"github.com/Lighty0410/ekadashi-server/pkg/storage/mongo"

	grpcServer "github.com/Lighty0410/ekadashi-server/pkg/server/grpc"
)

func main() {
	listner, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println("cannot connect to gRPC: ", err)
	}
	mongoService, err := mongo.NewService("localhost:27017")
	if err != nil {
		log.Fatalf("Could not create mongo service: %v", err)
	}
	newController := controller.NewController(mongoService)
	createController := grpcServer.CreateServer(newController)
	///
	service := grpcServer.ServerServer(createController)
	server := grpc.NewServer()
	grpcServer.RegisterServerServer(server, service)
	err = server.Serve(listner)
	if err != nil {
		log.Println("cannot listen gRPC server: ", err)
	}
}
