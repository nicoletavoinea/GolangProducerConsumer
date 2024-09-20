package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type StatusCode int

const (
	RECEIVED   StatusCode = iota //0
	PROCESSING                   //1
	DONE                         //2
)

type task struct {
	taskId             int32 //to make it static
	taskType           int8
	taskValue          int8
	taskState          StatusCode
	taskCreationTime   int64
	taskLastUpdateTime int64
}

func generateRandomTask() task {
	now := time.Now()

	var randomTask task
	randomTask.taskId = rand.Int32()
	randomTask.taskType = int8(rand.IntN(9))
	randomTask.taskValue = int8(rand.IntN(99))
	randomTask.taskState = RECEIVED
	randomTask.taskCreationTime = now.Unix() //secons elapsed since 1970
	randomTask.taskLastUpdateTime = randomTask.taskCreationTime
	return randomTask
}

func main() {
	fmt.Println("Initial commit for Producer service")

	fmt.Println(generateRandomTask())
	fmt.Println(generateRandomTask())
	fmt.Println(generateRandomTask())
}
