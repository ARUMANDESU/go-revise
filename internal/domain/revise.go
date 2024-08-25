package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

// ReviseItem is represents revision entity
type ReviseItem struct {
	ID            uuid.UUID
	AuthorID      uuid.UUID
	Name          string
	Description   string
	Tags          []string
	Iteration     ReviseIteration
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastRevisedAt time.Time
	NextReviseAt  time.Time
}
