package main

import (
	"flag"
	"log"

	"github.com/chienaeae/todo-go-grpc/client"
	"github.com/chienaeae/todo-go-grpc/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func testCreateTodo(todoClient *client.TodoClient) {
	todoClient.CreateTodo(sample.NewTodo())
}

func testGetTodos(todoClient *client.TodoClient) {
	todoClient.GetTodos()
}

func testUploadImage(todoClient *client.TodoClient) {
	todo := sample.NewTodo()
	todoClient.CreateTodo(todo)
	todoClient.UploadImage(todo.Id, "tmp/todo.png")
	todoClient.UploadImage(todo.Id, "tmp/aoi.jpeg")
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()

	log.Printf("dial server %s", *serverAddress)

	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.NewClient(*serverAddress, transportOption)
	if err != nil {
		log.Fatal("cannot dial server", err)
	}
	todoClient := client.NewTodoClient(cc)

	testCreateTodo(todoClient)
	testGetTodos(todoClient)
	testUploadImage(todoClient)
}
