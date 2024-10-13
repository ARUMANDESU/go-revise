package reviseitem

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
)

type SqliteRepo struct {
	db *sql.DB
}

// Save saves a revise item.
func (SqliteRepo) Save(ctx context.Context, item Aggregate) (_ error) {
	panic("not implemented") // TODO: Implement
}

// Update updates a revise item.
func (SqliteRepo) Update(ctx context.Context, id uuid.UUID, fn UpdateFn) (_ error) {
	panic("not implemented") // TODO: Implement
}
