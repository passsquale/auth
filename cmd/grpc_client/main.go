package main

import (
	"context"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50051"
	userID  = 12
)

func main() {
	cli, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer cli.Close()

	c := desc.NewUserV1Client(cli)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	log.Printf("User info: \n  %+v", resp.GetUser())

}
