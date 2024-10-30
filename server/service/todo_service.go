package service

import (
	"context"
	pb "todo-grpc/proto"
	"sync"
	"time"
	"github.com/google/uuid"
)

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
	mu    sync.Mutex
	tasks map[string]*pb.Task
}

func NewTodoServer() *TodoServer {
	return &TodoServer{
		tasks: make(map[string]*pb.Task),
	}
}

func (s *TodoServer) CreateTask(ctx context.Context, req *pb.Task) (*pb.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	task := &pb.Task{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   time.Now().Unix(),
	}
	
	s.tasks[id] = task
	return task, nil
}

func (s *TodoServer) ListTasks(_ *pb.Empty, stream pb.TodoService_ListTasksServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if err := stream.Send(task); err != nil {
			return err
		}
	}
	return nil
}