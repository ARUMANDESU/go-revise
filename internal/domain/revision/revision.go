package revision

import (
	"time"

	"github.com/gofrs/uuid"
)

type Revision struct {
	ID           uuid.UUID
	ReviseItemID uuid.UUID
	RevisedAt    time.Time
	//Notes        string // maybe in the future
}
