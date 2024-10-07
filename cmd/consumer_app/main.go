package main

import (
	"log"
	"net"

	"github.com/nicoletavoinea/GolangProducerConsumer/api/handler"
	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task" // proto package
	"google.golang.org/grpc"
)

func main() {
	// Set up a listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the TaskService server
	proto.RegisterTaskServiceServer(grpcServer, &handler.TaskServiceServer{})

	// Start serving
	log.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
