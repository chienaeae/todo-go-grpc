package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/chienaeae/todo-go-grpc/pb"
	"github.com/chienaeae/todo-go-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()

	todoStore := service.NewInMemoryTodoStore()
	imageStore := service.NewDiskImageStore("img")
	feedbackStore := service.NewInMemoryFeedbackStore()
	todoServer := service.NewTodoServer(
		todoStore,
		imageStore,
		feedbackStore,
	)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTodoServiceServer(srv, todoServer)
	reflection.Register(srv)

	log.Printf("Start GRPC server at %s", listener.Addr().String())
	err = srv.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
