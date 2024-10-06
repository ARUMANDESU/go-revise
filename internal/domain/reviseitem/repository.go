package reviseitem

import (
	"context"

	"github.com/gofrs/uuid"
)

type Repository interface {
	// Save saves a revise item.
	Save(ctx context.Context, item ReviseItem) error
	// GetById retrieves a revise item by ID.
	GetById(ctx context.Context, id uuid.UUID) (ReviseItem, error)
	// GetByUserID retrieves a list of revise items by user ID.
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]ReviseItem, error)
	// Update updates a revise item.
	Update(ctx context.Context, item ReviseItem) error
	// Delete deletes a revise item.
	Delete(ctx context.Context, id uuid.UUID) error
}
