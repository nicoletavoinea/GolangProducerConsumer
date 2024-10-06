package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicoletavoinea/GolangProducerConsumer/api/handler"
	"github.com/nicoletavoinea/GolangProducerConsumer/internal/database"
	"github.com/nicoletavoinea/GolangProducerConsumer/internal/metrics"
)

func main() {

	db := database.OpenDatabase()             //open database
	metrics.CreatePrometheusMetricsGeneral()  //initialize prometheus metrics
	go metrics.StartPrometheusServer(":2112") //start prometheus server

	//configure http setup
	router := mux.NewRouter()
	router.HandleFunc("/task", handler.HandleTask).Methods("POST")

	log.Println("Starting the Task Consumer on :8081...")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	database.CloseDB(db) //close database

}
