package handler

import (
	"context"
	"log"

	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task"
	"github.com/nicoletavoinea/GolangProducerConsumer/internal/definitions"
)

type TaskServiceServer struct {
	proto.UnimplementedTaskServiceServer // Required for forward compatibility
}

// HandleTask is the implementation of the HandleTask RPC method
func (s *TaskServiceServer) HandleTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	task := req.GetTask()

	// Log or process the task
	log.Printf("Received Task: %v", task)

	// Simulate task processing
	task.TaskState = proto.TaskState_PROCESSING

	// Return a response with the updated task state
	return &proto.TaskResponse{
		Task: task,
	}, nil
}

func (s *TaskServiceServer) ProcessAndSendTask(ctx context.Context, req *proto.Empty) (*proto.TaskResponse, error) {
	// Simulate creating a new task
	task := definitions.GenerateRandomTask()

	// Log task creation
	log.Println("Created Task:", task)

	// Simulate sending task to consumer
	task.TaskState = proto.TaskState_DONE

	// Return a response with the updated task state
	return &proto.TaskResponse{
		Task: task,
	}, nil
}

func SendTask(client proto.TaskServiceClient) {
	task := definitions.GenerateRandomTask()

	//response, err := client.HandleTask(context.Background(), &proto.TaskRequest{
	_, err := client.HandleTask(context.Background(), &proto.TaskRequest{
		Task: task,
	})

	if err != nil {
		log.Printf("Failed to call HandleTask: %v", err)
		return
	}
	log.Printf("Sent: %v", task)
	//log.Printf("Response from server: %v", response.GetTask())
}
