package reviseitem

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
)

var (
	ErrInvalidUserID         = fmt.Errorf("invalid userID")
	ErrNameRequired          = fmt.Errorf("name is required")
	ErrDescriptionRequired   = fmt.Errorf("description is required")
	ErrInvalidTags           = fmt.Errorf("tags are invalid")
	ErrInvalidNextRevisionAt = fmt.Errorf("nextRevisionAt is required")
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
		return ReviseItem{}, ErrInvalidUserID
	}
	if err := validateName(name); err != nil {
		return ReviseItem{}, err
	}
	if err := validateDescription(description); err != nil {
		return ReviseItem{}, err
	}
	if err := validateTags(tags); err != nil {
		return ReviseItem{}, err
	}
	if err := validateNextRevisionAt(nextRevisionAt); err != nil {
		return ReviseItem{}, err
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
	if err := validateName(name); err != nil {
		return err
	}

	r.name = name
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) UpdateDescription(description string) error {
	if err := validateDescription(description); err != nil {
		return err
	}

	r.description = description
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) UpdateTags(tags valueobject.StringArray) error {
	if err := validateTags(tags); err != nil {
		return err
	}

	r.tags = tags
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) PurgeTags() {
	r.tags = nil
	r.updatedAt = time.Now()
}

func (r *ReviseItem) UpdateNextRevisionAt(nextRevisionAt time.Time) error {
	if err := validateNextRevisionAt(nextRevisionAt); err != nil {
		return err
	}

	r.nextRevisionAt = nextRevisionAt
	r.lastRevisedAt = time.Now()
	r.updatedAt = time.Now()

	return nil
}
