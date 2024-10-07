package definitions

import (
	"math/rand"
	"time"

	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task"
)

type StatusCode int

const (
	RECEIVED   StatusCode = iota //0
	PROCESSING                   //1
	DONE                         //2
)

type Task struct {
	TaskId             int32      `json:"id"`
	TaskType           int32      `json:"type"`
	TaskValue          int32      `json:"value"`
	TaskState          StatusCode `json:"state"`
	TaskCreationTime   int64      `json:"creationtime"`
	TaskLastUpdateTime int64      `json:"lastupdatetime"`
}

func GenerateRandomTaskPrev() Task { //generate random values for type & value fields
	return Task{
		TaskType:  int32(rand.Intn(10)),
		TaskValue: int32(rand.Intn(100)),
		TaskState: RECEIVED,
	}
}

func GenerateRandomTask() *proto.Task {
	taskID := rand.Int31n(1000)
	taskType := rand.Int31n(5)
	taskValue := rand.Int31n(500)
	creationTime := time.Now().Unix()
	lastUpdateTime := time.Now().Unix()

	return &proto.Task{
		TaskId:             taskID,
		TaskType:           taskType,
		TaskValue:          taskValue,
		TaskState:          proto.TaskState_RECEIVED,
		TaskCreationTime:   creationTime,
		TaskLastUpdateTime: lastUpdateTime,
	}
}
