package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aditya43/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("From client..")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	client := greetpb.NewGreetServiceClient(clientConn)

	doUnary(client)
	doServerStreaming(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aditya",
			LastName:  "Hajare",
		},
	}
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error response received: %v", err)
	}
	fmt.Println(res.Result)
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aditya",
			LastName:  "Hajare",
		},
	}
	resStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error response received: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err != io.EOF {
			// We've reached the end of stream
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		res := msg.GetResult()
		fmt.Printf("Response: %v\n", res)
	}
}
