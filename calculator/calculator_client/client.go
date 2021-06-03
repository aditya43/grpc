package main

import (
	"context"
	"fmt"
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
