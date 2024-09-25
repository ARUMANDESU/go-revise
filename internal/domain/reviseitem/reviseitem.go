package reviseitem

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
)

type ReviseItem struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Name        string
	Description string
	Tags        valueobject.StringArray

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	NextRevisionAt time.Time
	LastRevisedAt  time.Time
}
