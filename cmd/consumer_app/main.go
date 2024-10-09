package main

import (
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"
	"github.com/nicoletavoinea/GolangProducerConsumer/api/handler"
	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task" // proto package
	database "github.com/nicoletavoinea/GolangProducerConsumer/internal/database"
	"github.com/nicoletavoinea/GolangProducerConsumer/internal/definitions"
	metrics "github.com/nicoletavoinea/GolangProducerConsumer/internal/metrics"
	"google.golang.org/grpc"
)

var (
	version = "0.0.0"
)

func main() {
	definitions.Version_handling(version)

	//open database
	myDatabase := database.OpenDatabase()

	//start prometheus
	metrics.CreatePrometheusMetricsGeneral()  //initialize prometheus metrics
	go metrics.StartPrometheusServer(":2112") //start prometheus server

	time.Sleep(time.Duration(10 * time.Second))

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

	database.CloseDB(myDatabase)
}
