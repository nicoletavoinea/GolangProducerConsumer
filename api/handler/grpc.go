package handler

import (
	"context"
	"log"
	"time"

	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task"
	database "github.com/nicoletavoinea/GolangProducerConsumer/internal/database"
	sqlc "github.com/nicoletavoinea/GolangProducerConsumer/internal/database/sqlc"
	definitions "github.com/nicoletavoinea/GolangProducerConsumer/internal/definitions"
	metrics "github.com/nicoletavoinea/GolangProducerConsumer/internal/metrics"
	zl "github.com/rs/zerolog/log"
)

type TaskServiceServer struct {
	proto.UnimplementedTaskServiceServer // Required for forward compatibility
}

func (s *TaskServiceServer) HandleTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	task := req.GetTask()
	//log.Printf("Received: %v", task)
	zl.Info().Int32("TaskId", task.TaskId).Int32("TaskType", task.TaskType).Int32("TaskValue", task.TaskValue).Str("State", string(task.TaskState)).Msg("Received by consumer")

	//Update ststus in database to PROCESSING
	database.UpdateTaskState(task.TaskId, proto.TaskState_PROCESSING)

	//update metrics
	metrics.IncreaseProcessedTasks()

	// Simulate task processing by sleeping for a duration based on task value
	time.Sleep(time.Duration(task.TaskValue) * time.Second)

	//Update ststus in database to DONE
	database.UpdateTaskState(task.TaskId, proto.TaskState_DONE)

	//update metrics
	metrics.IncreaseDoneTasks()

	// Return a response with the updated task state
	return &proto.TaskResponse{
		Task: task,
	}, nil
}

func AddToDatabaseAndSendTask(client proto.TaskServiceClient, queries *sqlc.Queries) {
	//generate random task
	task := definitions.GenerateRandomTask()

	//increase metrics
	metrics.IncreaseTotalTasksAndValue(task.TaskType, task.TaskValue)

	//add task to database
	database.AddTaskToDatabase(task, queries)

	//send task to consumer via grpc
	SendTask(client, task)
}

func SendTask(client proto.TaskServiceClient, task *proto.Task) {
	_, err := client.HandleTask(context.Background(), &proto.TaskRequest{
		Task: task,
	})

	if err != nil {
		log.Printf("Failed to call HandleTask: %v", err)
		return
	}
	//log.Printf("Sent: %v", task)
	zl.Info().Int32("TaskId", task.TaskId).Int32("TaskType", task.TaskType).Int32("TaskValue", task.TaskValue).Str("State", string(sqlc.TaskStateRECEIVED)).Msg("Sent by producer")
}
