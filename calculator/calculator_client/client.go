package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aditya43/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	doClientStreaming(client)
	doBidirectionalStreaming(client)
	doErrorHandlingExampleUnary(client)
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

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error opening client stream: %v", err)
	}

	numbers := []int32{2, 4, 5, 6, 21, 27, 45, 88, 92}

	for _, num := range numbers {
		fmt.Printf("Sending number: %v\n", num)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: num,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v\n", err)
	}
	fmt.Printf("The average is: %v", res.GetAverage())
}

func doBidirectionalStreaming(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error opening client stream: %v", err)
	}

	waitChan := make(chan struct{})
	go func() {
		defer stream.CloseSend()
		numbers := []int32{7, 2, 9, 24, 1, 3, 99, 55}
		for _, num := range numbers {
			fmt.Printf("Sending number: %v\n", num)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: num,
			})
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		defer close(waitChan)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving response: %v\n", err)
				break
			}
			fmt.Printf("New Maximum: %v\n", res.GetMaximum())
		}
	}()

	<-waitChan
}

func doErrorHandlingExampleUnary(client calculatorpb.CalculatorServiceClient) {
	numbers := []int32{12, -4}

	for _, num := range numbers {
		res, err := client.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
			Number: num,
		})

		if err != nil {
			respErr, ok := status.FromError(err)
			if ok {
				// Actual error from gRPC (use error)
				fmt.Println(respErr.Details()...)
				fmt.Println(respErr.Message(), respErr.Code())
				return
			} else {
				// Internal/Framework error
				log.Fatalf("Error calling SquareRoot: %v", err)
				return
			}
		}

		fmt.Printf("SquareRoot of number %v is: %v", num, res.GetNumberRoot())
	}
}
