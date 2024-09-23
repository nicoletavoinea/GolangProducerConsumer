package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand/v2"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	database "github.com/nicoletavoinea/GolangProducerConsumer/database/sqlc"
	//	"github.com/prometheus/client_golang/prometheus"
	//	"github.com/prometheus/client_golang/prometheus/promauto"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
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

func generateRandomTask() task {
	//now := time.Now()

	var randomTask task
	randomTask.TaskType = int8(rand.IntN(9))
	randomTask.TaskValue = int8(rand.IntN(99))
	randomTask.TaskState = RECEIVED
	//randomTask.TaskCreationTime = now.Unix() //secons elapsed since 1970
	//randomTask.TaskLastUpdateTime = randomTask.TaskCreationTime
	return randomTask
}

func processAndSend(taskToSend task, queries *database.Queries) {
	taskToSend, err := addToDatabase(taskToSend, queries)
	if err != nil {
		log.Printf("Error inserting task into database: %v\n", err)
	}
	sendTask(taskToSend)
}

func addToDatabase(taskToAdd task, queries *database.Queries) (task, error) {

	taskData, err := queries.AddTask(context.Background(), database.AddTaskParams{
		Param1: int64(taskToAdd.TaskType),
		Param2: int64(taskToAdd.TaskValue),
	})

	if err != nil {
		log.Printf("Error inserting task: %v\n", err)
		return taskToAdd, err
	}
	taskToAdd.TaskId = int32(taskData.ID)
	taskToAdd.TaskCreationTime = taskData.Creationtime
	taskToAdd.TaskLastUpdateTime = taskToAdd.TaskCreationTime
	return taskToAdd, nil
}

func sendTask(taskToSend task) {
	jsonTask, err := json.Marshal(taskToSend)
	if err != nil {
		log.Fatalf("Failed during Marshal function: %v", err)
	}

	response, err := http.Post("http://localhost:8081/task", "application/json", bytes.NewBuffer(jsonTask))
	if err != nil {
		log.Fatalf("Failed to post task: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Failed to send task to consumer: received status code %d", response.StatusCode)
	} else {
		log.Printf("Task sent to consumer: %+v", taskToSend)
	}
}

func runSchema(db *sql.DB, schemaFilePath string) error {
	schema, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	return err
}

func main() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := runSchema(db, "../database/schema.sql"); err != nil {
		log.Fatalf("Failed to run schema: %v", err)
	}

	queries := database.New(db)

	//create & send task
	tosend := generateRandomTask()
	processAndSend(tosend, queries)

}
