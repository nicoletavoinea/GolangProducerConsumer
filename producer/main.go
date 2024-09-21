package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
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
	now := time.Now()

	var randomTask task
	randomTask.TaskId = rand.Int32()
	randomTask.TaskType = int8(rand.IntN(9))
	randomTask.TaskValue = int8(rand.IntN(99))
	randomTask.TaskState = RECEIVED
	randomTask.TaskCreationTime = now.Unix() //secons elapsed since 1970
	randomTask.TaskLastUpdateTime = randomTask.TaskCreationTime
	return randomTask
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

func main() {
	//http.Handle("/metrics", promhttp.Handler())
	//http.ListenAndServe(":2112", nil)

	fmt.Println(generateRandomTask())
	fmt.Println(generateRandomTask())
	fmt.Println(generateRandomTask())
	tosend := generateRandomTask()
	sendTask(tosend)
	tosend = generateRandomTask()
	sendTask(tosend)
	tosend = generateRandomTask()
	sendTask(tosend)
}
