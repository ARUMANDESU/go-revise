package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/storage"
)

func (s *Storage) GetRevise(ctx context.Context, id string) (domain.ReviseItem, error) {
	const op = "storage.sqlite.revise.getRevise"

	query := `
		SELECT id, user_id, name, description, tags, iteration, created_at, updated_at, last_rivised_at, next_revision_at
		FROM revise_items
		WHERE id = ?
		`
	row := s.DB.QueryRowContext(ctx, query, id)

	var revise domain.ReviseItem
	err := row.Scan(
		&revise.ID,
		&revise.UserID,
		&revise.Name,
		&revise.Description,
		&revise.Tags,
		&revise.Iteration,
		&revise.CreatedAt,
		&revise.UpdatedAt,
		&revise.LastRevisedAt,
		&revise.NextRevisionAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ReviseItem{}, domain.WrapErrorWithOp(storage.ErrNotFound, op, "failed to get revise")
		}
		return domain.ReviseItem{}, domain.WrapErrorWithOp(err, op, "failed to get revise")
	}

	return revise, nil
}

func (s *Storage) ListRevises(ctx context.Context, dto domain.ListReviseItemDTO) ([]domain.ReviseItem, domain.PaginationMetadata, error) {
	const op = "storage.sqlite.revise.listRevises"

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, user_id, name, description, tags, iteration, created_at, updated_at, last_rivised_at, next_revision_at
		FROM revise_items
		WHERE user_id = ?
		ORDER BY %s %s
		LIMIT ? OFFSET ?`,
		dto.Sort.Field, dto.Sort.Order)

	rows, err := s.DB.QueryContext(ctx, query, dto.UserID, dto.Pagination.Limit(), dto.Pagination.Offset())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.PaginationMetadata{}, domain.WrapErrorWithOp(storage.ErrNotFound, op, "failed query to list revises")
		}
		return nil, domain.PaginationMetadata{}, domain.WrapErrorWithOp(err, op, "failed query to list revises")
	}
	defer rows.Close()

	var (
		total   int
		revises []domain.ReviseItem
	)
	for rows.Next() {
		var revise domain.ReviseItem
		err := rows.Scan(
			&total,
			&revise.ID, &revise.UserID, &revise.Name, &revise.Description, &revise.Tags, &revise.Iteration,
			&revise.CreatedAt, &revise.UpdatedAt, &revise.LastRevisedAt, &revise.NextRevisionAt,
		)
		if err != nil {
			return nil, domain.PaginationMetadata{}, domain.WrapErrorWithOp(err, op, "failed to list revises")
		}

		revises = append(revises, revise)
	}

	return revises, domain.CalculatePaginationMetadata(total, dto.Pagination.Page, dto.Pagination.PageSize), nil

}

func (s *Storage) CreateRevise(ctx context.Context, revise domain.ReviseItem) error {
	const op = "storage.sqlite.revise.createRevise"

	query := `
		INSERT INTO revise_items (id, user_id, name, description, tags, iteration, created_at, updated_at, last_rivised_at, next_revision_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
	args := []interface{}{
		revise.ID,
		revise.UserID,
		revise.Name,
		revise.Description,
		revise.Tags.Value(), // sqlite does not support array type, so we need to convert it to string
		revise.Iteration,
		revise.CreatedAt,
		revise.UpdatedAt,
		revise.LastRevisedAt,
		revise.NextRevisionAt,
	}

	_, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to create revise")
	}

	return nil
}

func (s *Storage) UpdateRevise(ctx context.Context, revise domain.ReviseItem) error {
	const op = "storage.sqlite.revise.updateRevise"

	query := `UPDATE revise_items
		SET name = ?, description = ?, tags = ?, iteration = ?, updated_at = ?, last_rivised_at = ?, next_revision_at = ?
		WHERE id = ?`

	args := []interface{}{
		revise.Name, revise.Description, revise.Tags.Value(), revise.Iteration,
		revise.UpdatedAt, revise.LastRevisedAt, revise.NextRevisionAt, revise.ID,
	}

	result, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to update revise")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return domain.WrapErrorWithOp(storage.ErrNotFound, op, "failed to update revise")
	}

	return nil
}

func (s *Storage) DeleteRevise(ctx context.Context, id string) error {
	const op = "storage.sqlite.revise.deleteRevise"

	query := `DELETE FROM revise_items WHERE id = ?`

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to delete revise")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return domain.WrapErrorWithOp(storage.ErrNotFound, op, "failed to delete revise")
	}

	return nil
}

func (s *Storage) GetScheduledItems(ctx context.Context) ([]domain.ScheduledItem, error) {
	const op = "storage.sqlite.revise.getScheduledItems"

	query := `
		SELECT ri.id, u.id, u.telegram_id, ri.name, ri.description, ri.tags, ri.iteration, ri.created_at, ri.updated_at, ri.last_rivised_at, ri.next_revision_at 
		FROM revise_items ri JOIN users u ON ri.user_id = u.id
		WHERE ri.next_revision_at <= datetime('now','+3 Hour');
	`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, domain.WrapErrorWithOp(err, op, "failed query to get scheduled items")
	}
	defer rows.Close()

	var scheduledItems []domain.ScheduledItem
	for rows.Next() {
		var scheduled domain.ScheduledItem
		err := rows.Scan(
			&scheduled.ReviseItem.ID, &scheduled.User.ID, &scheduled.User.TelegramID,
			&scheduled.ReviseItem.Name, &scheduled.ReviseItem.Description, &scheduled.ReviseItem.Tags, &scheduled.ReviseItem.Iteration,
			&scheduled.ReviseItem.CreatedAt, &scheduled.ReviseItem.UpdatedAt, &scheduled.ReviseItem.LastRevisedAt, &scheduled.ReviseItem.NextRevisionAt,
		)
		if err != nil {
			return nil, domain.WrapErrorWithOp(err, op, "failed to get scheduled items")
		}

		scheduledItems = append(scheduledItems, scheduled)
	}

	return scheduledItems, nil
}
