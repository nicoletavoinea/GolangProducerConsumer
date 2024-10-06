package definitions

import (
	"math/rand"
)

type StatusCode int

const (
	RECEIVED   StatusCode = iota //0
	PROCESSING                   //1
	DONE                         //2
)

type Task struct {
	TaskId             int32      `json:"id"`
	TaskType           int8       `json:"type"`
	TaskValue          int8       `json:"value"`
	TaskState          StatusCode `json:"state"`
	TaskCreationTime   int64      `json:"creationtime"`
	TaskLastUpdateTime int64      `json:"lastupdatetime"`
}

func GenerateRandomTask() Task { //generate random values for type & value fields
	return Task{
		TaskType:  int8(rand.Intn(10)),
		TaskValue: int8(rand.Intn(100)),
		TaskState: RECEIVED,
	}
}
