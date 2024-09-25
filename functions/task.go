package functions

import (
	"context"
	"log"
	"math/rand"

	database "github.com/nicoletavoinea/GolangProducerConsumer/database/sqlc"
)

type StatusCode int

const (
	RECEIVED   StatusCode = iota //0
	PROCESSING                   //1
	DONE                         //2
)

type task struct {
	TaskId             int32      `json:"id"`
	TaskType           int8       `json:"type"`
	TaskValue          int8       `json:"value"`
	TaskState          StatusCode `json:"state"`
	TaskCreationTime   int64      `json:"creationtime"`
	TaskLastUpdateTime int64      `json:"lastupdatetime"`
}

func GenerateRandomTask() task {
	return task{
		TaskType:  int8(rand.Intn(9)),
		TaskValue: int8(rand.Intn(99)),
		TaskState: RECEIVED,
	}
}

func ProcessAndSendTask(task task, db *database.Queries) error {
	// Add task to database
	updatedTask, err := AddTaskToDatabase(task, db)
	if err != nil {
		return err
	}
	log.Println("Task added to db:", updatedTask)

	// Send task to consumer
	err = SendTaskToConsumer(updatedTask)
	if err != nil {
		return err
	}
	log.Println("Task sent to consumer:", updatedTask)
	return nil
}

func AddTaskToDatabase(task task, db *database.Queries) (task, error) {
	taskData, err := db.AddTask(context.Background(), database.AddTaskParams{
		Param1: int64(task.TaskType),
		Param2: int64(task.TaskValue),
	})
	if err != nil {
		log.Printf("Error inserting task: %v\n", err)
		return task, err
	}
	task.TaskId = int32(taskData.ID)
	task.TaskCreationTime = taskData.Creationtime
	task.TaskLastUpdateTime = task.TaskCreationTime
	return task, nil
}

func UpdateTaskState(taskID int32, status StatusCode) (database.Task, error) {
	state := ""
	if status == PROCESSING {
		state = "PROCESSING"
	} else if status == DONE {
		state = "DONE"
	}

	updatedTask, err := Queries.UpdateTask(context.Background(), database.UpdateTaskParams{
		Param1: int64(taskID),
		Param2: state,
	})
	if err != nil {
		log.Printf("Error updating task: %v\n", err)
		return updatedTask, err
	}
	return updatedTask, nil
}
