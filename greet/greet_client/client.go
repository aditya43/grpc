package main

import (
	"fmt"
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
	fmt.Printf("Created client: %v", client)
}
