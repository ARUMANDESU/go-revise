package reviseitem

import (
	"context"

	"github.com/gofrs/uuid"
)

type UpdateFn func(item *Aggregate) (*Aggregate, error)

type Repository interface {
	// Save saves a revise item.
	Save(ctx context.Context, item Aggregate) error
	// GetById retrieves a revise item by ID.
	GetById(ctx context.Context, id uuid.UUID) (Aggregate, error)
	// GetByUserID retrieves a list of revise items by user ID.
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Aggregate, error)
	// Update updates a revise item.
	Update(ctx context.Context, id uuid.UUID, fn UpdateFn) error
}
