package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicoletavoinea/GolangProducerConsumer/functions"
)

func main() {
	db := functions.OpenDatabase()
	functions.CreatePrometheusMetricsGeneral()
	go functions.StartPrometheusServer(":2112")

	time.Sleep(time.Duration(10) * time.Second)

	router := mux.NewRouter()
	router.HandleFunc("/task", functions.HandleTask).Methods("POST")

	log.Println("Starting the Task Consumer on :8081...")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	functions.CloseDB(db)

}
