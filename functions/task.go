package functions

import (
	"log"
	"math/rand"
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
		TaskType:  int8(rand.Intn(10)),
		TaskValue: int8(rand.Intn(100)),
		TaskState: RECEIVED,
	}
}

func ProcessAndSendTask() (task, error) {
	task := GenerateRandomTask()

	// Add task to database
	updatedTask, err := AddTaskToDatabase(task, Queries)
	if err != nil {
		return updatedTask, err
	}
	log.Println("Task added to db:", updatedTask)

	//update metrics
	IncreaseTotalTasksAndValue(task.TaskType, task.TaskValue)

	// Send task to consumer
	err = SendTaskToConsumer(updatedTask)
	if err != nil {
		return updatedTask, err
	}
	log.Println("Task sent to consumer:", updatedTask)
	return updatedTask, nil
}
