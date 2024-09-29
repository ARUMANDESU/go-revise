package reviseitem

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
)

// ReviseItem represents a revise item.
// A revise item is a thing that needs to be revised.
// It is mutable.
type ReviseItem struct {
	id     uuid.UUID
	userID uuid.UUID

	name        string
	description string
	tags        valueobject.StringArray

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time

	nextRevisionAt time.Time
	lastRevisedAt  time.Time
}

// NewReviseItemID creates a new revise item ID.
func NewReviseItemID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

// NewReviseItem creates a new revise item. It returns an error if the arguments are invalid.
func NewReviseItem(userID uuid.UUID, name, description string, tags valueobject.StringArray, nextRevisionAt time.Time) (ReviseItem, error) {
	if userID == uuid.Nil {
		return ReviseItem{}, fmt.Errorf("userID is required")
	}
	if name == "" {
		return ReviseItem{}, fmt.Errorf("name is required")
	}
	if description == "" {
		return ReviseItem{}, fmt.Errorf("description is required")
	}
	if !tags.IsValid() {
		return ReviseItem{}, fmt.Errorf("tags are invalid")
	}
	if nextRevisionAt.IsZero() {
		return ReviseItem{}, fmt.Errorf("nextRevisionAt is required")
	}

	return ReviseItem{
		id:             NewReviseItemID(),
		userID:         userID,
		name:           name,
		description:    description,
		tags:           tags,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
		nextRevisionAt: nextRevisionAt,
	}, nil
}

func (r *ReviseItem) UpdateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	r.name = name
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) UpdateDescription(description string) error {
	if description == "" {
		return fmt.Errorf("description is required")
	}

	r.description = description
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) UpdateTags(tags valueobject.StringArray) error {
	if !tags.IsValid() {
		return fmt.Errorf("tags are invalid")
	}

	r.tags = tags
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) UpdateNextRevisionAt(nextRevisionAt time.Time) error {
	if nextRevisionAt.IsZero() {
		return fmt.Errorf("nextRevisionAt is required")
	}

	r.nextRevisionAt = nextRevisionAt
	r.lastRevisedAt = time.Now()
	r.updatedAt = time.Now()

	return nil
}
