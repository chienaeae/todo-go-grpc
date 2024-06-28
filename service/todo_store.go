package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/chienaeae/todo-go-grpc/pb"
	"github.com/jinzhu/copier"
)


var ErrAlreadyExists = errors.New("record already exists")

type TodoStore interface {
	Save(todo *pb.Todo) error
	GetById(id string) (*pb.Todo, error)
	GetMany(ctx context.Context, found func(todo *pb.Todo) error) error
}

type InMemoryTodoStore struct {
	mutex sync.RWMutex
	data map[string]*pb.Todo
}

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{
		data: make(map[string]*pb.Todo),
	}
}

func (store *InMemoryTodoStore) Save(todo *pb.Todo) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[todo.Id] != nil {
		return ErrAlreadyExists
	}

	other, err := deepCopy(todo)
	if err != nil {
		return err
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryTodoStore) GetById(id string) (*pb.Todo, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	todo := store.data[id]
	if todo == nil {
		return nil, nil
	}

	return deepCopy(todo)
}

func (store *InMemoryTodoStore) GetMany(ctx context.Context, found func(todo *pb.Todo) error ) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, todo := range store.data {
		err := ctx.Err()
		if err == context.Canceled || err == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return nil
		}

		other, err := deepCopy(todo)
		if err != nil {
			return err
		}

		err = found(other)
		if err != nil {
			return err
		}
	}
	return nil
}

func deepCopy(todo *pb.Todo) (*pb.Todo, error) {
	other := &pb.Todo{}

	if err := copier.Copy(other, todo); err != nil {
		return nil, fmt.Errorf("cannot copy todo data: %w", err)
	}

	return other, nil

}