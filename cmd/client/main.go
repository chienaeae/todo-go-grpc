package main

import (
	"flag"
	"log"
	"time"

	"github.com/chienaeae/todo-go-grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	username        = "admin"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func dial(serverAddress string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	log.Printf("dial server %s", serverAddress)
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	allOpts := append([]grpc.DialOption{transportOption}, opts...)
	return grpc.NewClient(serverAddress, allOpts...)
}

func testAuth(cc *grpc.ClientConn) {
	authClient := client.NewAuthClient(cc, username, password)
	testLogin(authClient)
}

func testTodo(cc *grpc.ClientConn) {
	todoClient := client.NewTodoClient(cc)
	testCreateTodo(todoClient)
	testGetTodos(todoClient)
	testUploadImage(todoClient)
	testCreateFeedbacks(todoClient)
}

func authMethods() map[string]bool {
	const todoServicePath = "/todoGoGrpc.TodoService/"
	return map[string]bool{
		todoServicePath + "CreateTodo":   true,
		todoServicePath + "GetTodos":     true,
		todoServicePath + "GetTodo":      true,
		todoServicePath + "FeedbackTodo": true,
		todoServicePath + "UploadImage":  true,
	}
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	service := flag.String("service", "todo", "execute service target")
	flag.Parse()

	cc1, err := dial(*serverAddress)
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	if *service == "todo" {
		authClient := client.NewAuthClient(cc1, username, password)
		interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
		if err != nil {
			log.Fatal("cannot create auth interceptor: ", err)
		}

		cc2, err := dial(
			*serverAddress,
			grpc.WithUnaryInterceptor(interceptor.Unary()),
			grpc.WithStreamInterceptor(interceptor.Stream()),
		)
		if err != nil {

		}

		testTodo(cc2)
	} else if *service == "auth" {
		testAuth(cc1)
	} else {
		log.Fatalf("unknown service: %s", *service)
	}

}
