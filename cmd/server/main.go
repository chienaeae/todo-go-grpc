package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/chienaeae/todo-go-grpc/pb"
	"github.com/chienaeae/todo-go-grpc/service"
	"google.golang.org/grpc"
)

func main () {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()

	todoStore := service.NewInMemoryTodoStore()
	todoServer := service.NewTodoServer(
		todoStore,
	)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTodoServiceServer(srv, todoServer)

	log.Printf("Start GRPC server at %s", listener.Addr().String())
	err = srv.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}