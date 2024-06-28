package sample

import (
	"math/rand"

	"github.com/chienaeae/todo-go-grpc/pb"
)

func NewTodo() *pb.Todo {
	length := rand.Intn(20);
	todo := &pb.Todo {
		Id: randomID(),
		Title: randomString(length),
	}

	return todo
}