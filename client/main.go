package main

import (
	"context"
	"io"
	"log"
	pb "todo-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)
	
	// Criar uma tarefa
	task, err := client.CreateTask(context.Background(), &pb.Task{
		Title:       "Aprender gRPC",
		Description: "Desenvolver um projeto MVP com gRPC",
	})
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	log.Printf("Task created: %v", task)

	// Listar tarefas
	stream, err := client.ListTasks(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("could not list tasks: %v", err)
	}

	for {
		task, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Task: %v", task)
	}
} 