package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/metadata"

	grpcServer "github.com/Lighty0410/ekadashi-server/pkg/server/grpc"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051", grpc.WithInsecure(),
	)
	if err != nil {
		log.Println("cannot connect to gRPC server: ", err)
	}

	ctx := context.Background()
	md := metadata.Pairs(
		"token", "hmmmetoya",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	client := grpcServer.NewServerClient(conn)
	u, err := client.HandleRegistration(ctx, &grpcServer.User{User: "asidhasi", Password: "hmmmetoya"})
	if err != nil {
		fmt.Println(err)
		return
	}

	sess, err := client.HandleLogin(ctx)
	fmt.Println(u.Response)
}
