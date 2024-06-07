package session

import "context"

// ErrSessionNotFound represents error where appropriate selenium is not found
type ErrSessionNotFound error

type Repository interface {
	FindAll() ([]*Session, error)
	FindByID(id string) (*Session, error)
	Delete(id string) error
	Save(s *Session) error
	Watch(ctx context.Context) <-chan *Session
	Quota() (int, int, error)
}
