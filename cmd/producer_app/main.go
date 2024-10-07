package main

import (
	"log"
	"sync"
	"time"

	"github.com/nicoletavoinea/GolangProducerConsumer/api/handler"
	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task" // Import your generated package
	"google.golang.org/grpc"
)

const msgRate = 3 //message sending rate (messages/second)
var wg sync.WaitGroup

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewTaskServiceClient(conn)

	//generate tasks for 75s
	timeout := time.After(75 * time.Second)

	// Create a ticker that triggers every 500 milliseconds
	ticker := time.NewTicker(1000 / msgRate * time.Millisecond)
	defer ticker.Stop() // Stop the ticker when the program ends

	// Simulate sending a message every 500 milliseconds
	for {
		select {
		case <-ticker.C: // Every time the ticker ticks, send a task
			// Create a new task to send to the server
			go handler.SendTask(client)

		case <-timeout: // Triggers after timeout elapses
			ticker.Stop() // Stop the ticker
			wg.Wait()
		}
	}
}
