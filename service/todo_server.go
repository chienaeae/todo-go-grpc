package service

import (
	"context"
	"errors"
	"log"

	"github.com/chienaeae/todo-go-grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
	todoStore TodoStore
}

func NewTodoServer(store *InMemoryTodoStore) *TodoServer {
	return &TodoServer{
		todoStore: store,
	}
}


func (server *TodoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	todo := req.GetTodo()
	if len(todo.Id) > 0 {
		if _, err := uuid.Parse(todo.Id); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "todo ID is not a valid UUID: %v", err)
		}
	}else{
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new todo ID: %v", err)
		}
		todo.Id = id.String();
	}

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	err := server.todoStore.Save(todo)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists){
			code = codes.AlreadyExists
			
		}
		return nil, status.Errorf(code, "cannot save todo to the store: %v", err)
	}

	log.Printf("saved todo with id: %s", todo.Id)
	
	res := &pb.CreateTodoResponse{
		Id: todo.Id,
	}

	return res, nil
}

func (server *TodoServer) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	id := req.GetId()
	todo, err := server.todoStore.GetById(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if todo == nil {
		return nil, status.Errorf(codes.NotFound, "cannot find todo with ID: %v", id)
	}

	res := &pb.GetTodoResponse{
		Todo: todo,
	}
	return res, nil
}

func (server *TodoServer) GetTodos(req *pb.GetTodosRequest, stream pb.TodoService_GetTodosServer) error {
	log.Printf("receiving todos stream")
	err := server.todoStore.GetMany(
		stream.Context(),
		func(todo *pb.Todo) error {
			res := &pb.GetTodosResponse {
				Todo: todo,
			}

			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent todo with id: %s", todo.GetId())
			return nil
		},
	)

	if err != nil { 
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is cancelled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil;
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return nil
}