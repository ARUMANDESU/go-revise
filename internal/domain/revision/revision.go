package revision

import (
	"time"

	"github.com/gofrs/uuid"
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
