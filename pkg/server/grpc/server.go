package grpc

import (
	"log"
	"net"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	api "github.com/Lighty0410/ekadashi-server/pkg/server/grpc/api"
	"google.golang.org/grpc"
)

// NewGrpcServer sets up a new TCP route and creates a new gRPC server.
func NewGrpcServer(c *controller.Controller) (*grpc.Server, error) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println("cannot connect to gRPC: ", err)
	}
	service := NewService(c)
	server := grpc.NewServer()
	api.RegisterEkadashiServer(server, service)
	go func() {
		err = server.Serve(listener)
		if err != nil {
			log.Printf("cannot listen to gRPC server: %v", err)
		}
	}()
	return server, nil
}
