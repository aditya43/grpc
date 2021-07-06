package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aditya43/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	doClientStreaming(client)
	doBiDirectionalStreaming(client)

	doGreetWithDeadlineExample(client, 1) // Should timeout
	doGreetWithDeadlineExample(client, 5) // Should complete
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

func doClientStreaming(client greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Aditya",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Avantika",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nishigandha",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jane",
			},
		},
	}

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving client streaming response: %v", err)
	}
	fmt.Printf("LongGreet response: %v", res)
}

func doBiDirectionalStreaming(client greetpb.GreetServiceClient) {
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error opening client stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Aditya",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Avantika",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nishigandha",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jane",
			},
		},
	}

	waitChan := make(chan struct{})

	// Send requests
	go func(stream greetpb.GreetService_GreetEveryoneClient) {
		defer stream.CloseSend()
		for _, req := range requests {
			fmt.Printf("Sending message: %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
	}(stream)

	// Receive responses
	go func() {
		defer close(waitChan)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving server stream response: %v", err)
				break
			}
			fmt.Printf("Received Response: %v\n", res.GetResult())
		}
	}()

	<-waitChan
}

func doGreetWithDeadlineExample(client greetpb.GreetServiceClient, timeoutDuration time.Duration) {
	req := &greetpb.GreetDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aditya",
			LastName:  "Hajare",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration*time.Second)
	defer cancel()

	res, err := client.GreetDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout (5 seconds) reached! Deadline exceeded!")
			} else {
				fmt.Printf("Unexpected error received: %v", statusErr)
			}
		} else {
			log.Fatalf("Error calling GreetDeadline rpc: %v", err)
		}
		return
	}
	fmt.Println(res.Result)
}
