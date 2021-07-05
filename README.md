# gRPC Using Golang
My personal notes, projects and best practices.

## Author
Aditya Hajare ([Linkedin](https://in.linkedin.com/in/aditya-hajare)).

## Current Status
WIP (Work In Progress)!

## License
Open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT).

-----------

## What is gRPC?
- gRPC is a free and open source framework developed by Google, Square and other companies.
- gRPC is part of the Cloud Native Computation Foundation (CNCF). - like Docker and Kubernetes for example.
- At a high level, gRPC allows us to define REQUEST and RESPONSE for RPC (Remote Procedure Calls) and handles all the rest for us.
- gRPC is modern, fast and efficient, built on top of HTTP/2, low latency, supports streaming, language independent, and makes it super easy to plug in authentication, load balancing, logging and monitoring.
- RPC is not a new concept (CORBA had this before).
- With gRPC, RPC is implemented very cleanly and solves a lot of problems.
- At the core of gRPC, we need to define messages and services using `Protocol Buffers`.
- The rest of the gRPC code will be generated for us and we will have to provide an implementation for it.
- One `.proto` file works for over 12 programming languages (server and client), and allows us to use a framework that scales to millions of RPC per second.

-----------

## Why Protocol Buffers over JSON:
- gRPC uses Protocol Buffers for communication.
- Payload size comparison: Protocol Buffers vs. JSON:
```json
// 55 bytes
{
    "age": 35,
    "first_name": "Aditya",
    "last_name": "Hajare"
}
```
```proto
// 20 bytes
message Person {
    int32 age = 1;
    string first_name = 2;
    string last_name = 3;
}
```
- Looking at above comparison, we save in network bandwidth.
- Parsing JSON is actually CPU intensive (because the format is Human Readable).
- Parsing Protocol Buffers (Binary Format) is less CPU intensive because it's closer to how machine represents data.

-----------

## HTTP 2:
- gRPC leverages HTTP/2 as a backbone for communication.
- HTTP 1.1 opens a new TCP connection to a server for each request.
- HTTP 1.1 does not compress headers (Headers are plaintext).
- HTTP 1.1 only works with Request/Response mechanism (No server push).
- HTTP 1.1 was originally composed of 2 commands:
    * GET: to ask for content.
    * POST: to send content.
- HTTP 2 was released in 2015. It has been battle tested for many years. (And was before that tested by Google under the name SPDY).
- HTTP 2 supports `multiplexing`:
    * The client and server can push messages in parallel over the same TCP connection.
    * This will greatly reduce latency.
- HTTP 2 supports `server push`:
    * Server can push streams (multiple messages) for one request from the client.
    * This saves lot of round trips (latency).
- HTTP 2 supports `headers compression`.
- HTTP 2 is binary.
- HTTP 2 is secure (SSL is not required but recommended by default).

-----------

## Types of APIs in gRPC:
- **Unary API**: Classic Request/Response API.
- **Server Streaming API**: Client will send `one` message to the server and will receive `many` responses from the server, possibly an infinite number.
    * `Streaming Server API` are well suited for when the server needs to send a lot of data (big data).
    * We can use `Streaming Server API` when the server needs to `PUSH` data to the client without having client `REQUEST` for the more. For e.g. Live Feed, Chat etc..
- **Client Streaming API**: Client sends multiple chunks to server as in stream of requests and server responds with a single response.
    * `Streaming Client API` is well suited when the client needs to send lot of data (big data).
    * We can use `Streaming Client API` when the server processing is expensive and should happen as the client sends data.
    * It is also useful when client needs to **PUSH** data to the server without really expecting a response.
    * In gRPC, Client streaming calls are defined using the keywork `stream`.
- **Bi-Directional Streaming API**: Client and server both sends and receives request/response chunks in stream.
```proto
service GreetService {
    // Unary API
    rpc Greet(GreetRequest) returns (GreetResponse) {};

    // Streaming Server API
    rpc GreetManyTimes(GreetManyTimesRequest) returns (stream GreetManyTimesResponse) {};

    // Streaming Client API
    rpc LongGreet(stream LongGreetRequest) returns (LongGreetResponse) {};

    // Bi-Directional Streaming API
    rpc GreetEveryone(stream GreetEveryoneRequest) returns (stream GreetEveryoneResponse) {};
}
```

-----------

## Scalability in gRPC:
- gRPC servers are Asynchronous by default.
- That means they do not block threads on Request.
- Therefore each gRPC server can serve millions of requests in parallel.
- gRPC clients has a choice of being Asynchronous or Synchronous (blocking).
- gRPC clients can perform client side load balancing.
- Google has 10 Billion gRPC requests being made per second internally.