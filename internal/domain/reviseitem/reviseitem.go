package reviseitem

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

// ReviseItem represents a revise item.
// A revise item is a thing that needs to be revised.
// It is mutable.
type ReviseItem struct {
	id     uuid.UUID
	userID uuid.UUID

	name        string
	description string
	tags        valueobject.Tags

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

type NewReviseItemArgs struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Name           string
	Description    string
	Tags           valueobject.Tags
	NextRevisionAt time.Time
}

// NewReviseItem creates a new revise item. It returns an error if the arguments are invalid.
func NewReviseItem(args NewReviseItemArgs) (*ReviseItem, error) {
	if args.ID == uuid.Nil {
		return nil, ErrInvalidID
	}
	if args.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	args.Name = strings.TrimSpace(args.Name)
	if err := validateName(args.Name); err != nil {
		return nil, err
	}
	args.Description = strings.TrimSpace(args.Description)
	if err := validateDescription(args.Description); err != nil {
		return nil, err
	}
	if err := valueobject.ValidateTags(args.Tags); err != nil {
		return nil, err
	}
	if err := validateNextRevisionAt(args.NextRevisionAt); err != nil {
		return nil, err
	}

	now := time.Now()
	return &ReviseItem{
		id:             args.ID,
		userID:         args.UserID,
		name:           args.Name,
		description:    args.Description,
		tags:           args.Tags,
		createdAt:      now,
		updatedAt:      now,
		nextRevisionAt: args.NextRevisionAt,
	}, nil
}

func (r *ReviseItem) ID() uuid.UUID {
	return r.id
}

func (r *ReviseItem) UserID() uuid.UUID {
	return r.userID
}

func (r *ReviseItem) Name() string {
	return r.name
}

func (r *ReviseItem) Description() string {
	return r.description
}

func (r *ReviseItem) Tags() valueobject.Tags {
	return r.tags
}

func (r *ReviseItem) CreatedAt() time.Time {
	return r.createdAt
}

func (r *ReviseItem) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r *ReviseItem) DeletedAt() *time.Time {
	return r.deletedAt
}

func (r *ReviseItem) NextRevisionAt() time.Time {
	return r.nextRevisionAt
}

func (r *ReviseItem) LastRevisedAt() time.Time {
	return r.lastRevisedAt
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

func (r *ReviseItem) AddTags(tags valueobject.Tags) error {
	if !tags.IsValid() {
		return errs.NewIncorrectInputError("invalid tags", "invalid-tags")
	}
	if err := valueobject.ValidateTags(tags); err != nil {
		return err
	}

	r.tags.AddTags(tags)
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) RemoveTags(tags valueobject.Tags) error {
	if !tags.IsValid() {
		return errs.NewIncorrectInputError("invalid tags", "invalid-tags")
	}

	r.tags.RemoveTags(tags)
	r.updatedAt = time.Now()

	return nil
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

func (r *ReviseItem) MarkAsDeleted() {
	now := time.Now()
	r.deletedAt = &now
	r.updatedAt = now
}

func (r *ReviseItem) Restore() {
	r.deletedAt = nil
	r.updatedAt = time.Now()
}

func (r *ReviseItem) IsOwner(userID uuid.UUID) bool {
	return r.userID == userID
}

func (r *ReviseItem) CanModify(userID uuid.UUID) bool {
	return r.userID == userID
}
