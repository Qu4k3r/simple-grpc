# Detalhamento da Comunicação gRPC no Projeto Todo List

## 1. Visão Geral do gRPC

O gRPC é um framework de RPC (Remote Procedure Call) desenvolvido pelo Google que utiliza Protocol Buffers como linguagem de definição de interface (IDL) e formato de serialização de dados.

## 2. Definição do Serviço (Protocol Buffers)

### 2.1 Arquivo de Definição (todo.proto)
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

### 2.2 Explicação dos Componentes

- **Service TodoService**: Define dois métodos RPC:
  - `CreateTask`: Método unário (request-response simples)
  - `ListTasks`: Método de streaming do servidor (server-streaming)

- **Message Task**: Estrutura de dados principal com campos:
  - `id`: Identificador único da tarefa
  - `title`: Título da tarefa
  - `description`: Descrição detalhada
  - `completed`: Status de conclusão
  - `created_at`: Timestamp de criação

## 3. Fluxo de Comunicação

### 3.1 Criação de Tarefa (CreateTask)

1. **Cliente inicia a requisição**:
```go
task, err := client.CreateTask(context.Background(), &pb.Task{
    Title:       "Aprender gRPC",
    Description: "Desenvolver um projeto MVP com gRPC",
})
```

2. **Servidor processa a requisição**:
```go
func (s *TodoServer) CreateTask(ctx context.Context, req *pb.Task) (*pb.Task, error) {
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
```

### 3.2 Listagem de Tarefas (ListTasks)

1. **Cliente inicia o streaming**:
```go
stream, err := client.ListTasks(context.Background(), &pb.Empty{})
for {
    task, err := stream.Recv()
    if err == io.EOF {
        break
    }
    log.Printf("Task: %v", task)
}
```

2. **Servidor envia as tarefas**:
```go
func (s *TodoServer) ListTasks(_ *pb.Empty, stream pb.TodoService_ListTasksServer) error {
    for _, task := range s.tasks {
        if err := stream.Send(task); err != nil {
            return err
        }
    }
    return nil
}
```

## 4. Tipos de Comunicação gRPC Implementados

### 4.1 Unário (CreateTask)
- Uma única requisição do cliente
- Uma única resposta do servidor
- Síncrono
- Ideal para operações CRUD simples

### 4.2 Server Streaming (ListTasks)
- Uma única requisição do cliente
- Múltiplas respostas do servidor
- Assíncrono
- Ideal para buscar conjuntos de dados ou monitorar atualizações

## 5. Segurança e Concorrência

### 5.1 Mutex para Concorrência
```go
type TodoServer struct {
    mu    sync.Mutex
    tasks map[string]*pb.Task
}
```
- Protege o acesso concorrente ao mapa de tarefas
- Garante operações thread-safe

### 5.2 Contexto
```go
func (s *TodoServer) CreateTask(ctx context.Context, req *pb.Task)
```
- Permite cancelamento de operações
- Gerencia timeouts
- Propaga metadados entre cliente e servidor

## 6. Geração de Código

O protoc gera automaticamente:
- Interfaces do cliente
- Interfaces do servidor
- Estruturas de dados
- Métodos de serialização/deserialização

## 7. Benefícios da Implementação

1. **Performance**
   - Serialização binária eficiente
   - Conexões HTTP/2 multiplexadas
   - Streaming bidirecional

2. **Tipo-Seguro**
   - Contratos fortemente tipados
   - Verificação em tempo de compilação
   - Autocompletar IDE

3. **Escalabilidade**
   - Suporte a load balancing
   - Conexões persistentes
   - Streaming eficiente

## 8. Considerações de Produção

1. **Monitoramento**
   - Implementar interceptors para logging
   - Adicionar métricas de performance
   - Monitorar latência e erros

2. **Resiliência**
   - Implementar retry policies
   - Adicionar circuit breakers
   - Gerenciar timeouts adequadamente

3. **Segurança**
   - Adicionar TLS para criptografia
   - Implementar autenticação
   - Validar inputs

## 9. Próximos Passos Sugeridos

1. Implementar mais métodos gRPC:
   - Atualização de tarefas
   - Deleção de tarefas
   - Streaming bidirecional para atualizações em tempo real

2. Adicionar persistência de dados:
   - Integração com banco de dados
   - Cache distribuído

3. Melhorar a segurança:
   - Adicionar autenticação JWT
   - Implementar TLS
   - Adicionar validação de dados

4. Implementar observabilidade:
   - Tracing distribuído
   - Métricas de performance
   - Logging estruturado 