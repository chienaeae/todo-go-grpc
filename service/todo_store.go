package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

type TodoStore interface {
	Save(todo *Todo) error
	GetById(id string) (*Todo, error)
	GetMany(ctx context.Context, fromUser string, found func(todo *Todo) error) error
}

type Todo struct {
	ID       string
	Title    string
	FromUser string
}

type InMemoryTodoStore struct {
	mutex sync.RWMutex
	data  map[string]*Todo
}

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{
		data: make(map[string]*Todo),
	}
}

func (store *InMemoryTodoStore) Save(todo *Todo) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[todo.ID] != nil {
		return ErrAlreadyExists
	}

	other, err := deepCopy(todo)
	if err != nil {
		return err
	}

	store.data[other.ID] = other
	return nil
}

func (store *InMemoryTodoStore) GetById(id string) (*Todo, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	todo := store.data[id]
	if todo == nil {
		return nil, nil
	}

	return deepCopy(todo)
}

func (store *InMemoryTodoStore) GetMany(ctx context.Context, fromUser string, found func(todo *Todo) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, todo := range store.data {
		if todo.FromUser != fromUser {
			continue
		}

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

func deepCopy(todo *Todo) (*Todo, error) {
	other := &Todo{}

	if err := copier.Copy(other, todo); err != nil {
		return nil, fmt.Errorf("cannot copy todo data: %w", err)
	}

	return other, nil

}
