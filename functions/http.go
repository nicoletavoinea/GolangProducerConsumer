package functions

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HandleTask(w http.ResponseWriter, r *http.Request) {
	var receivedTask task

	// Parse the JSON body
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
	log.Printf("Received task: %v", receivedTask)

	// Update task state to PROCESSING
	updatedTask, err := UpdateTaskState(receivedTask.TaskId, PROCESSING)
	if err != nil {
		log.Println("Error updating task to PROCESSING :(\n", err)
	}
	//IncreaseProcessedTasks()
	IncreaseProcessedTasks()
	log.Println("Task updated to: ", updatedTask)

	//Sleep fot value miliseconds
	time.Sleep(time.Duration(receivedTask.TaskValue) * time.Second) //should be miliseconds

	// Update task state to DONE
	updatedTask, err = UpdateTaskState(receivedTask.TaskId, DONE)
	if err != nil {
		log.Println("Error updating task to DONE :(\n", err)
	}
	IncreaseDoneTasks()
	log.Println("Task updated to: ", updatedTask)

	//Ok status
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(receivedTask)
}

func SendTaskToConsumer(taskToSend task) error {
	jsonTask, err := json.Marshal(taskToSend)
	if err != nil {
		return err
	}

	response, err := http.Post("http://localhost:8081/task", "application/json", bytes.NewBuffer(jsonTask))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Failed to send task to consumer: received status code %d", response.StatusCode)
		return err
	}
	log.Printf("Task sent to consumer: %+v", taskToSend)
	return nil
}
