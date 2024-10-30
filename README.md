# Todo List gRPC com Docker

Este é um projeto de exemplo de uma Todo List usando gRPC e Docker, implementado em Go.

## Pré-requisitos

- Docker
- Docker Compose
- Make (opcional, mas recomendado)

## Estrutura do Projeto

```
todo-grpc/
├── docker-compose.yml
├── Dockerfile.server
├── Dockerfile.client
├── Dockerfile.protoc
├── Makefile
├── README.md
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

## Configuração Inicial

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd todo-grpc
```

2. Inicialize o módulo Go:
```bash
go mod init todo-grpc
go mod tidy
```

## Como Executar

### Usando Make

1. Gere os arquivos Protocol Buffers:
```bash
make proto
```

2. Construa as imagens Docker:
```bash
make build
```

3. Execute o servidor:
```bash
make run-server
```

4. Em outro terminal, execute o cliente:
```bash
make run-client
```

### Sem Make (usando Docker Compose diretamente)

1. Gere os arquivos Protocol Buffers:
```bash
docker-compose run --rm protoc
```

2. Construa as imagens:
```bash
docker-compose build
```

3. Execute o servidor:
```bash
docker-compose up server
```

4. Em outro terminal, execute o cliente:
```bash
docker-compose up client
```

## Funcionalidades

O projeto implementa um serviço gRPC simples com duas operações:

1. `CreateTask`: Cria uma nova tarefa
2. `ListTasks`: Lista todas as tarefas existentes (usando streaming)

## Limpeza

Para limpar os containers e arquivos gerados:

```bash
make clean
```

Ou usando Docker Compose diretamente:

```bash
docker-compose down
```

## Desenvolvimento

### Modificando o Protocol Buffer

1. Edite o arquivo `proto/todo.proto`
2. Regenere os arquivos Go:
```bash
make proto
```

### Hot Reload

O servidor está configurado com volumes Docker, então as alterações no código são refletidas automaticamente sem necessidade de reconstruir as imagens.

### Logs

Para ver os logs dos serviços:

```bash
# Logs do servidor
docker-compose logs -f server

# Logs do cliente
docker-compose logs -f client
```

## Troubleshooting

### Problemas Comuns

1. Se o cliente não conseguir conectar ao servidor:
   - Verifique se o servidor está rodando (`docker-compose ps`)
   - Verifique se a porta 50051 está disponível

2. Se os arquivos proto não forem gerados:
   - Verifique se o comando `make proto` foi executado
   - Verifique se há erros na sintaxe do arquivo .proto

3. Se as dependências não forem encontradas:
   - Execute `go mod tidy` para atualizar as dependências

### Portas Utilizadas

- gRPC Server: 50051

## Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Licença

Este projeto está sob a licença MIT. 