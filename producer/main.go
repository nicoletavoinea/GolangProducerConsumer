package main

import (
	"log"

	functions "github.com/nicoletavoinea/GolangProducerConsumer/functions"
)

func main() {

	//open database & get queries(global var)
	db := functions.OpenDatabase()

	log.Println("Icerc sa fac ceva")
	//create & send task
	tosend := functions.GenerateRandomTask()
	log.Println("Trying to send: %v", tosend)

	functions.ProcessAndSendTask(tosend, functions.Queries)

	functions.CloseDB(db)

}
