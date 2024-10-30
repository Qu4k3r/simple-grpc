# Guia de Desenvolvimento: Todo List com gRPC usando Docker

## 1. Estrutura do Projeto

```
todo-grpc/
├── docker-compose.yml
├── Dockerfile.server
├── Dockerfile.client
├── Dockerfile.protoc
├── Makefile
├── go.mod
├── proto/
│   └── todo.proto
├── server/
│   ├── main.go
│   └── service/
│       └── todo_service.go
└── client/
    └── main.go
```

## 2. Arquivos de Configuração

### 2.1 Dockerfile.protoc
```dockerfile
FROM golang:1.21-alpine

RUN apk add --no-cache protobuf-dev protoc

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /app
```

### 2.2 Dockerfile.server
```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server server/main.go

EXPOSE 50051

CMD ["/server"]
```

### 2.3 Dockerfile.client
```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /client client/main.go

CMD ["/client"]
```

### 2.4 docker-compose.yml
```yaml
version: '3.8'

services:
  protoc:
    build:
      context: .
      dockerfile: Dockerfile.protoc
    volumes:
      - .:/app
    command: ["sh", "-c", "protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/todo.proto"]

  server:
    build:
      context: .
      dockerfile: Dockerfile.server
    ports:
      - "50051:50051"
    volumes:
      - .:/app
    command: ["go", "run", "server/main.go"]

  client:
    build:
      context: .
      dockerfile: Dockerfile.client
    volumes:
      - .:/app
    depends_on:
      - server
    command: ["go", "run", "client/main.go"]
```

### 2.5 Makefile
```makefile
.PHONY: proto build run-server run-client clean

proto:
	docker-compose run --rm protoc

build:
	docker-compose build

run-server:
	docker-compose up server

run-client:
	docker-compose up client

clean:
	docker-compose down
	rm -f server/server client/client
```

### 2.6 go.mod
```bash
go mod init todo-grpc
go mod tidy
```

## 3. Arquivos do Projeto

### 3.1 proto/todo.proto
```protobuf
syntax = "proto3";

package todo;
option go_package = "todo-grpc/proto";

service TodoService {
    rpc CreateTask (Task) returns (Task) {}
    rpc ListTasks (Empty) returns (stream Task) {}
}

message Task {
    string id = 1;
    string title = 2;
    string description = 3;
    bool completed = 4;
    int64 created_at = 5;
}

message Empty {}
```

### 3.2 server/service/todo_service.go
```go
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
```

### 3.3 server/main.go
```go
package main

import (
    "log"
    "net"
    "todo-grpc/service"
    pb "todo-grpc/proto"
    "google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterTodoServiceServer(s, service.NewTodoServer())
    
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

### 3.4 client/main.go
```go
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
```

## 4. Passos para Desenvolvimento

1. **Inicializar o projeto:**
```bash
mkdir todo-grpc
cd todo-grpc
# Copiar todos os arquivos de configuração e código
```

2. **Gerar código protobuf:**
```bash
make proto
```

3. **Construir as imagens:**
```bash
make build
```

4. **Executar o servidor:**
```bash
make run-server
```

5. **Em outro terminal, executar o cliente:**
```bash
make run-client
```

## 5. Comandos Úteis

- **Recompilar protobuf:**
```bash
make proto
```

- **Reconstruir imagens:**
```bash
make build
```

- **Limpar ambiente:**
```bash
make clean
```

## 6. Desenvolvimento

1. O servidor está configurado para reiniciar automaticamente quando houver mudanças no código
2. Para testar novas funcionalidades, você pode modificar o código e o servidor será recarregado
3. Para testar o cliente com novas mudanças, execute `make run-client` novamente

## 7. Dicas de Desenvolvimento

1. Use `docker-compose logs -f server` para ver os logs do servidor
2. Use `docker-compose logs -f client` para ver os logs do cliente
3. Para adicionar dependências:
   - Adicione no go.mod
   - Execute `docker-compose build` para reconstruir as imagens
4. Para debugging, você pode adicionar mais logs usando `log.Printf()`
