package service

import (
	"context"

	"github.com/chienaeae/todo-go-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	jwtManager *JWTManager
	userStore  UserStore
}

func NewAuthServer(jwtManager *JWTManager, userStore UserStore) *AuthServer {
	return &AuthServer{
		jwtManager: jwtManager,
		userStore:  userStore,
	}
}

func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	if username == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password are required")
	}

	user, err := server.userStore.Find(username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect username/password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.LoginResponse{AccessToken: token}
	return res, nil
}
