package definitions

import (
	"math/rand"
	"time"

	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task"
)

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
