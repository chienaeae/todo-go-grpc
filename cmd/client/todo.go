package main

import (
	"fmt"
	"strings"

	"github.com/chienaeae/todo-go-grpc/client"
	"github.com/chienaeae/todo-go-grpc/sample"
)

func testCreateTodo(todoClient *client.TodoClient) {
	todoClient.CreateTodo(sample.NewTodo())
}

func testGetTodos(todoClient *client.TodoClient) {
	todoClient.GetTodos()
}

func testUploadImage(todoClient *client.TodoClient) {
	todo := sample.NewTodo()
	todoClient.CreateTodo(todo)
	todoClient.UploadImage(todo.Id, "tmp/todo.png")
	todoClient.UploadImage(todo.Id, "tmp/aoi.jpeg")
}

func testCreateFeedbacks(todoClient *client.TodoClient) {
	todo := sample.NewTodo()
	todoClient.CreateTodo(todo)

	n := 3
	createFeedbacks := make([]client.CreateFeedback, n)
	for {
		for i := 0; i < n; i++ {
			createFeedbacks[i] = client.CreateFeedback{
				TodoID:  todo.Id,
				Content: sample.NewContent(),
			}
		}

		todoClient.FeedbackTodo(createFeedbacks)

		fmt.Print("continue to create feedback (y/N)?")
		var ans string
		fmt.Scan(&ans)
		if strings.ToLower(ans) != "y" {
			break
		}
	}
}
