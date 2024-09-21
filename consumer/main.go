package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type StatusCode int

const (
	RECEIVED   StatusCode = iota //0
	PROCESSING                   //1
	DONE                         //2
)

type task struct {
	TaskId             int32      `json:"id"` //to make it static
	TaskType           int8       `json:"type"`
	TaskValue          int8       `json:"value"`
	TaskState          StatusCode `json:"state"`
	TaskCreationTime   int64      `json:"creationtime"`
	TaskLastUpdateTime int64      `json:"lastupdatetime"`
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	var receivedTask task

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &receivedTask)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received task: %+v", receivedTask)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(receivedTask)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/task", handleTask).Methods("POST")

	log.Println("Starting the Task Consumer on :8081...")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
