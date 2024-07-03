package main

import (
	"flag"
	"log"

	"github.com/chienaeae/todo-go-grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func dial(serverAddress string) (*grpc.ClientConn, error) {
	log.Printf("dial server %s", serverAddress)
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	return grpc.NewClient(serverAddress, transportOption)
}

func testAuth(cc *grpc.ClientConn) {
	authClient := client.NewAuthClient(cc)
	testLogin(authClient)
}

func testTodo(cc *grpc.ClientConn) {
	todoClient := client.NewTodoClient(cc)
	testCreateTodo(todoClient)
	testGetTodos(todoClient)
	testUploadImage(todoClient)
	testCreateFeedbacks(todoClient)
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	service := flag.String("service", "todo", "execute service target")
	flag.Parse()

	cc, err := dial(*serverAddress)
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	if *service == "todo" {
		testTodo(cc)
	} else if *service == "auth" {
		testAuth(cc)
	} else {
		log.Fatalf("unknown service: %s", *service)
	}

}
