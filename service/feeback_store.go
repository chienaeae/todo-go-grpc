package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type FeedbackStore interface {
	Add(todoID string, content string) (*Feedback, error)
}

type Feedback struct {
	ID      string
	Content string
}

type InMemoryFeedbackStore struct {
	mutex     sync.RWMutex
	feedbacks map[string][]*Feedback
}

func NewInMemoryFeedbackStore() *InMemoryFeedbackStore {
	return &InMemoryFeedbackStore{
		feedbacks: make(map[string][]*Feedback),
	}
}

func (store *InMemoryFeedbackStore) Add(todoID string, content string) (*Feedback, error) {
	feedbackID, err := uuid.NewRandom()
	newFeedback := &Feedback{
		ID:      feedbackID.String(),
		Content: content,
	}

	if err != nil {
		return nil, fmt.Errorf("cannot generate feedback id: %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	feedbacks := store.feedbacks[todoID]
	if feedbacks == nil {
		feedbacks = make([]*Feedback, 0)
	}
	feedbacks = append(feedbacks, newFeedback)

	return newFeedback, nil
}
