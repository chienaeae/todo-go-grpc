package client

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/chienaeae/todo-go-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoClient struct {
	service pb.TodoServiceClient
}

func NewTodoClient(cc *grpc.ClientConn) *TodoClient {
	service := pb.NewTodoServiceClient(cc)
	return &TodoClient{service}
}

func (todoClient *TodoClient) CreateTodo (todo *pb.Todo) {
	log.Println("=== CreateTodo ===")

	req := &pb.CreateTodoRequest {
		Todo: todo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := todoClient.service.CreateTodo(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch (st.Code()) {
			case codes.AlreadyExists:
				log.Print("todo already exists")
			}
		}else {
			log.Fatal("cannot create todo: ", err)
		}
		return 
	}

	log.Printf("created todo with ID: %s", res.Id)
}

func (todoClient *TodoClient) GetTodos () {
	log.Println("=== GetTodos ===")
	req := &pb.GetTodosRequest {}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := todoClient.service.GetTodos(ctx, req)
	if err != nil {
		log.Fatal("cannot get todos: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}
		todo := res.GetTodo()
		log.Printf("<%s>", todo.Id)
		log.Print("title: ", todo.Title)
	}
}
