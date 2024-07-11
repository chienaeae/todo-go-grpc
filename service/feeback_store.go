package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type FeedbackStore interface {
	Add(todoID string, feedback *Feedback) (*Feedback, error)
	Find(todoID string) ([]*Feedback, error)
}

type Feedback struct {
	ID       string
	Content  string
	FromUser string
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

func (store *InMemoryFeedbackStore) Add(todoID string, feedback *Feedback) (*Feedback, error) {
	feedbackID, err := uuid.NewRandom()
	newFeedback := &Feedback{
		ID:       feedbackID.String(),
		Content:  feedback.Content,
		FromUser: feedback.FromUser,
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

	store.feedbacks[todoID] = feedbacks
	return deepCopyFeedback(newFeedback)
}

func (store *InMemoryFeedbackStore) Find(todoID string) ([]*Feedback, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	fs := store.feedbacks[todoID]
	return fs, nil
}

func deepCopyFeedback(feedback *Feedback) (*Feedback, error) {
	other := &Feedback{}

	if err := copier.Copy(other, feedback); err != nil {
		return nil, fmt.Errorf("cannot copy todo data: %w", err)
	}

	return other, nil

}
