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

	//ticker that starts routines having msgRate/second
	ticker := time.NewTicker(time.Second / msgRate)
	defer ticker.Stop()

	//start generating & sending tasks
	for range ticker.C {
		// Generate and send the task
		//log.Println("Goroutine started")
		go functions.ProcessAndSendTask()
	}

	//wait for all goroutines to finish
	wg.Wait()

	//close database
	functions.CloseDB(db)

}
