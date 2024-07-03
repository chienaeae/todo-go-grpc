package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/chienaeae/todo-go-grpc/pb"
	"github.com/chienaeae/todo-go-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func seedUsers(userStore service.UserStore) error {
	_, err := createUser(userStore, "admin", "secret", "admin")
	if err != nil {
		return err
	}

	return nil
}

func createUser(userStore service.UserStore, username, password, role string) (*service.User, error) {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return nil, err
	}

	err = userStore.Save(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}
	todoStore := service.NewInMemoryTodoStore()
	imageStore := service.NewDiskImageStore("img")
	feedbackStore := service.NewInMemoryFeedbackStore()
	todoServer := service.NewTodoServer(
		todoStore,
		imageStore,
		feedbackStore,
	)
	authServer := service.NewAuthServer(jwtManager, userStore)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTodoServiceServer(srv, todoServer)
	pb.RegisterAuthServiceServer(srv, authServer)
	reflection.Register(srv)

	log.Printf("Start GRPC server at %s", listener.Addr().String())
	err = srv.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
