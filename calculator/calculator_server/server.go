package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aditya43/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Received sum RPC request: %v", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()

	res := &calculatorpb.SumResponse{
		SumResult: firstNumber + secondNumber,
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
	calculatorpb.RegisterCalculatorServiceServer(grpcServer, &server{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
