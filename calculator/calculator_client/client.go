package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aditya43/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("From client..")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	client := calculatorpb.NewCalculatorServiceClient(clientConn)

	doUnary(client)
	doServerStreaming(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  2,
		SecondNumber: 3,
	}
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error response received: %v", err)
	}
	fmt.Println(res)
}

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 24,
	}
	stream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error response received: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something went wrong: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}
