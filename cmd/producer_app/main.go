package main

import (
	//"log"
	"sync"
	"time"

	"github.com/nicoletavoinea/GolangProducerConsumer/api/handler"
	"github.com/nicoletavoinea/GolangProducerConsumer/internal/database"
	"github.com/nicoletavoinea/GolangProducerConsumer/internal/metrics"
)

const msgRate = 3 //message sending rate (messages/second)

var wg sync.WaitGroup

func main() {
	//open database & get queries(global var)
	db := database.OpenDatabase()
	metrics.CreatePrometheusMetricsTypes()
	go metrics.StartPrometheusServer(":2113")

	//10 s before starts generating tasks
	time.Sleep(time.Duration(10) * time.Second)

	//generate tasks for 75s
	timeout := time.After(75 * time.Second)

	//ticker that starts routines having msgRate/second
	ticker := time.NewTicker(time.Second / msgRate)
	defer ticker.Stop()

	//start generating & sending tasks
	for {
		select {
		case <-ticker.C: // Triggers at msgRate per second
			// Start a new goroutine for processing the task
			go func() {
				handler.ProcessAndSendTask()
			}()

		case <-timeout: // Triggers after timeout elapses
			ticker.Stop()        // Stop the ticker
			wg.Wait()            //wait for all goroutines to finish
			database.CloseDB(db) //close database

		}
	}

}
