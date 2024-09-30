package main

import (
	//"log"

	"sync"
	"time"

	functions "github.com/nicoletavoinea/GolangProducerConsumer/functions"
)

const msgRate = 5

var wg sync.WaitGroup

func main() {
	//open database & get queries(global var)
	db := functions.OpenDatabase()
	functions.CreatePrometheusMetricsTypes()
	go functions.StartPrometheusServer(":2113")

	time.Sleep(time.Duration(10) * time.Second)

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
				functions.ProcessAndSendTask()
			}()

		case <-timeout: // Triggers after 1 minute
			//wait for all goroutines to finish
			ticker.Stop() // Stop the ticker
			wg.Wait()
			functions.CloseDB(db)

		}
	}

	//close database

}
