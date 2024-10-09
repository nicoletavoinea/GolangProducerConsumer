package main

import (
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/nicoletavoinea/GolangProducerConsumer/api/handler"
	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task" // Import your generated package
	database "github.com/nicoletavoinea/GolangProducerConsumer/internal/database"
	metrics "github.com/nicoletavoinea/GolangProducerConsumer/internal/metrics"

	"google.golang.org/grpc"
)

const msgRate = 3 //message sending rate (messages/second)
var wg sync.WaitGroup

func main() {
	//open database
	myDatabase := database.OpenDatabase()

	//start prometheus
	metrics.CreatePrometheusMetricsTypes()    //initialize prometheus metrics
	go metrics.StartPrometheusServer(":2113") //start prometheus server

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
			go handler.AddToDatabaseAndSendTask(client, database.Queries)

		case <-timeout: // Triggers after timeout elapses
			ticker.Stop() // Stop the ticker
			wg.Wait()
			database.CloseDB(myDatabase)
		}
	}
}
