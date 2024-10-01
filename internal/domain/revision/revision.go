package revision

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/multierr"
)

var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrInvalidReviseItemID = errors.New("invalid revise item id")
)

// Revision represents a revision of a revise item.
// A revision is a snapshot of a revise item at a certain time.
// It is used to track the history of a revise item.
// It is immutable.
type Revision struct {
	ID           uuid.UUID
	ReviseItemID uuid.UUID
	RevisedAt    time.Time
	//Notes        string // maybe in the future
}

func NewRevisionID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func NewRevision(reviseItemID uuid.UUID) (*Revision, error) {
	if reviseItemID == uuid.Nil {
		return nil, multierr.Combine(ErrInvalidArgument, ErrInvalidReviseItemID)
	}

	return &Revision{
		ID:           NewRevisionID(),
		ReviseItemID: reviseItemID,
		RevisedAt:    time.Now(),
	}, nil
}
