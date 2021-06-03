package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/aditya43/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

// Unary API
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

// Streaming Server API
func (s *server) GreetManyTimes(req *greetpb.GreetRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function invoked with request: %v", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		str := "Hello " + firstName + " " + lastName + " " + strconv.Itoa(i)
		res := &greetpb.GreetResponse{
			Result: str,
		}

		_ = stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
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
