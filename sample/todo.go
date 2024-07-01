package sample

import (
	"math/rand"

	"github.com/chienaeae/todo-go-grpc/pb"
)

func NewTodo() *pb.Todo {
	length := rand.Intn(20)
	todo := &pb.Todo{
		Id:    randomID(),
		Title: randomString(length),
	}

	return todo
}

func NewContent() string {
	wordsCount := randomInt(10, 20)
	content := ""

	for i := 0; i < wordsCount; i++ {
		wordLen := randomInt(3, 9)
		content += randomString(wordLen)
		if i != wordsCount-1 {
			content += " "
		}
	}

	return content
}
