package client

import (
	"context"
	"log"
	"time"

	"github.com/chienaeae/todo-go-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AuthClient struct {
	service pb.AuthServiceClient
}

func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	return &AuthClient{service: pb.NewAuthServiceClient(cc)}
}

func (client *AuthClient) Login(username, password string) {
	log.Println("=== Login ===")

	req := &pb.LoginRequest{
		Username: username,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := client.service.Login(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatalf("cannot login: %s", st.Message())
	}

	log.Printf("access token: %s", res.AccessToken)
}
