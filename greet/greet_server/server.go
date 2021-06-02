package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aditya43/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function invoked with request: %v", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	str := "Hello " + firstName + " " + lastName

	res := &greetpb.GreetResponse{
		Result: str,
	}

	return res, nil
}

func main() {
	fmt.Println("From server..")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
