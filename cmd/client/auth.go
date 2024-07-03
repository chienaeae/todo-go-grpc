package main

import (
	"github.com/chienaeae/todo-go-grpc/client"
)

func testLogin(authClient *client.AuthClient) {
	authClient.Login("admin", "secret")
}
