package main

import (
	"log"

	"github.com/chienaeae/todo-go-grpc/client"
)

func testLogin(authClient *client.AuthClient) {
	accessToken, err := authClient.Login()
	if err != nil {
		log.Fatalf("cannot login: %s", err)
	}

	log.Printf("access token: %s", accessToken)
}
