package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

func (todoClient *TodoClient) CreateTodo(todo *pb.Todo) {
	log.Println("=== CreateTodo ===")

	req := &pb.CreateTodoRequest{
		Todo: todo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := todoClient.service.CreateTodo(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				log.Print("todo already exists")
			}
		} else {
			log.Fatal("cannot create todo: ", err)
		}
		return
	}

	log.Printf("created todo with ID: %s", res.Id)
}

func (todoClient *TodoClient) GetTodos() {
	log.Println("=== GetTodos ===")
	req := &pb.GetTodosRequest{}

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

func (laptopClient *TodoClient) UploadImage(todoID string, imagePath string) {
	log.Println("=== UploadImage ===")
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_ImageInfo{
			ImageInfo: &pb.ImageInfo{
				TodoId:    todoID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

type CreateFeedback struct {
	TodoID  string
	Content string
}

func (laptopClient *TodoClient) FeebackTodo(createFeedbacks []CreateFeedback) {
	log.Println("=== FeebackTodo ===")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.FeebackTodo(ctx)
	if err != nil {
		log.Fatal("cannot feedback todo: ", err)
	}

	waitResponse := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more response")
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
			}

			log.Print("received response: ", res)
		}
	}()

	for _, createFeedback := range createFeedbacks {
		req := &pb.FeedbackTodoRequest{
			TodoId:  createFeedback.TodoID,
			Content: createFeedback.Content,
		}

		err := stream.Send(req)
		if err != nil {
			log.Fatalf("cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
		}

		log.Print("sent request: ", req)
	}

	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("cannot close send: %v", err)
	}

	err = <-waitResponse
	if err != nil {
		log.Fatalf("received error: %v", err)
	}
}
