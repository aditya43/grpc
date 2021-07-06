package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/aditya43/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Streaming Client API
func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function invoked with stream request")
	var result string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function invoked with stream request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := &greetpb.GreetEveryoneResponse{
			Result: "Hello " + firstName + "! ",
		}

		if err = stream.Send(result); err != nil {
			return err
		}
	}
}

// Unary API
func (s *server) GreetDeadline(ctx context.Context, req *greetpb.GreetDeadlineRequest) (*greetpb.GreetDeadlineResponse, error) {
	fmt.Printf("Greet function invoked with request: %v", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// The client cancelled the request
			fmt.Println("The client cancelled the request")
			return nil, status.Error(codes.Canceled, "Client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	str := "Hello " + firstName + " " + lastName

	res := &greetpb.GreetDeadlineResponse{
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
