package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/chienaeae/todo-go-grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// maxImagesize is 1 Megabyte
const maxImageSize = 1 << 20

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
	todoStore  TodoStore
	imageStore ImageStore
}

func NewTodoServer(todoStore TodoStore, imageStore ImageStore) *TodoServer {
	return &TodoServer{
		todoStore:  todoStore,
		imageStore: imageStore,
	}
}

func (server *TodoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	todo := req.GetTodo()
	if len(todo.Id) > 0 {
		if _, err := uuid.Parse(todo.Id); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "todo ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new todo ID: %v", err)
		}
		todo.Id = id.String()
	}

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	err := server.todoStore.Save(todo)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
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
			res := &pb.GetTodosResponse{
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

func (server *TodoServer) UploadImage(stream pb.TodoService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	todoID := req.GetImageInfo().GetTodoId()
	imageType := req.GetImageInfo().GetImageType()

	todo, err := server.todoStore.GetById(todoID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot find todo: %v", err))
	}

	if todo == nil {
		return logError(status.Errorf(codes.InvalidArgument, "todo id %s doesn't exist", todoID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		err = contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err = stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("receivd a chunk with size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	imageID, err := server.imageStore.Save(todoID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("saved image with id: %s, size: %d", imageID, imageSize)
	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is cancelled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return nil
}
