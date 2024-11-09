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
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Tags        valueobject.Tags
}

// NewReviseItem creates a new revise item. It returns an error if the arguments are invalid.
func NewReviseItem(args NewReviseItemArgs) (*ReviseItem, error) {
	op := errs.Op("domain.reviseitem.new_revise_item")
	if args.ID == uuid.Nil {
		return nil, errs.
			NewUnknownError(op, errs.ErrInvalidInput, "revise item id is nil").
			WithContext("args.id", args.ID)
	}
	if args.UserID == uuid.Nil {
		return nil, errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "user id must be provided").
			WithMessages([]errs.Message{{Key: "message", Value: "user id must be provided"}}).
			WithContext("args.UserID", args.UserID)
	}
	args.Name = strings.TrimSpace(args.Name)
	if err := validateName(args.Name); err != nil {
		return nil, errs.WithOp(op, err, "validating revise item name failed")
	}
	args.Description = strings.TrimSpace(args.Description)
	if err := validateDescription(args.Description); err != nil {
		return nil, errs.WithOp(op, err, "validating revise item description failed")
	}
	if err := valueobject.ValidateTags(args.Tags); err != nil {
		return nil, errs.WithOp(op, err, "validating revise item tags failed")
	}
	// if err := validateNextRevisionAt(args.NextRevisionAt); err != nil {
	// 	return nil, errs.WithOp(op, err, "validating, revise item next revision at failed")
	// }

	now := time.Now()
	return &ReviseItem{
		id:             args.ID,
		userID:         args.UserID,
		name:           args.Name,
		description:    args.Description,
		tags:           args.Tags,
		createdAt:      now,
		updatedAt:      now,
		nextRevisionAt: valueobject.DefaultReviewIntervals().Next(0),
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
	op := errs.Op("domain.reviseitem.update_name")
	if err := validateName(name); err != nil {
		return errs.WithOp(op, err, "name validation failed")
	}

	r.name = name
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) UpdateDescription(description string) error {
	op := errs.Op("domain.reviseitem.update_description")
	if err := validateDescription(description); err != nil {
		return errs.WithOp(op, err, "description validation failed")
	}

	r.description = description
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) AddTags(tags valueobject.Tags) error {
	op := errs.Op("domain.reviseitem.add_tags")
	if tags.IsEmpty() {
		return errs.NewIncorrectInputError(op, nil, "tags is empty").WithContext("tags", tags)
	}
	if err := valueobject.ValidateTags(tags); err != nil {
		return errs.WithOp(op, err, "tags validation failed")
	}

	r.tags.AddTags(tags)
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) RemoveTags(tags valueobject.Tags) error {
	op := errs.Op("domain.reviseitem.remove_tags")
	if tags.IsEmpty() {
		return errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "tags are empty").
			WithMessages([]errs.Message{{Key: "message", Value: "tags must be provided, got empty"}})
	}

	r.tags.RemoveTags(tags)
	r.updatedAt = time.Now()

	return nil
}

func (r *ReviseItem) UpdateNextRevisionAt(nextRevisionAt time.Time) error {
	op := errs.Op("domain.reviseitem.update_next_revision_at")
	if err := validateNextRevisionAt(nextRevisionAt); err != nil {
		return errs.WithOp(op, err, "next revision at validation failed")
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
