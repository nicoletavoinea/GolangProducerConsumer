package definitions

import (
	"math/rand"

	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task"
)

func GenerateRandomTask() *proto.Task {
	taskType := rand.Int31n(10)
	taskValue := rand.Int31n(100)

	return &proto.Task{
		TaskType:  taskType,
		TaskValue: taskValue,
		TaskState: proto.TaskState_RECEIVED,
	}
}
