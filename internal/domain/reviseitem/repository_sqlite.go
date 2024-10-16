package reviseitem

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/db/sqlc"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/pointers"
)

type SqliteRepo struct {
	db *sql.DB
}

// Save saves a revise item.
func (r *SqliteRepo) Save(ctx context.Context, item Aggregate) (_ error) {
	tags := item.Tags()
	args := sqlc.SaveReviseItemParams{
		ID:             item.id.String(),
		UserID:         item.userID.String(),
		Name:           item.name,
		Description:    sql.NullString{String: item.description, Valid: item.description != ""},
		Tags:           stringArrToString(tags.StringArray()),
		CreatedAt:      item.createdAt,
		UpdatedAt:      item.updatedAt,
		LastRevisedAt:  item.lastRevisedAt,
		NextRevisionAt: item.nextRevisionAt,
	}

	q := sqlc.New(r.db)

	return q.SaveReviseItem(ctx, args)
}

// Update updates a revise item.
func (r *SqliteRepo) Update(ctx context.Context, id uuid.UUID, fn UpdateFn) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	qtx := sqlc.New(tx)

	reviseItemModel, err := qtx.GetReviseItem(ctx, id.String())
	if err != nil {
		return err
	}

	reviseItem, err := modelToReviseItem(reviseItemModel)
	if err != nil {
		return err
	}

	aggregate := NewAggregate(&reviseItem)

	aggregate, err = fn(aggregate)
	if err != nil {
		return err
	}

	for _, r := range aggregate.Revisions() {
		args := sqlc.CreateRevisionParams{
			ID:           r.ID().String(),
			ReviseItemID: aggregate.ID().String(),
			RevisedAt:    r.RevisedAt(),
		}

		err = qtx.CreateRevision(ctx, args)
		if err != nil {
			return err
		}
	}

	tags := aggregate.Tags()
	return qtx.UpdateReviceItem(ctx, sqlc.UpdateReviceItemParams{
		Name:           aggregate.Name(),
		Description:    sql.NullString{String: aggregate.Description(), Valid: aggregate.Description() != ""},
		Tags:           stringArrToString(tags.StringArray()),
		CreatedAt:      aggregate.CreatedAt(),
		UpdatedAt:      aggregate.UpdatedAt(),
		LastRevisedAt:  aggregate.LastRevisedAt(),
		NextRevisionAt: aggregate.NextRevisionAt(),
		ID:             aggregate.ID().String(),
	})
}

func stringArrToString(arr []string) sql.NullString {
	// transform the array into a string: ["a","b","c"] -> "a,b,c"
	if arr == nil || len(arr) == 0 {
		return sql.NullString{String: "", Valid: false}
	}

	// convert to string
	stringValue := strings.Join(arr, ",")
	if len(stringValue) == 0 {
		return sql.NullString{String: "", Valid: false}
	}

	return sql.NullString{String: stringValue, Valid: true}
}

func stringToStringArr(str sql.NullString) []string {
	if !str.Valid {
		return nil
	}

	return strings.Split(str.String, ",")
}

func modelToReviseItem(model sqlc.ReviseItem) (ReviseItem, error) {
	return ReviseItem{
		id:             uuid.FromStringOrNil(model.ID),
		userID:         uuid.FromStringOrNil(model.UserID),
		name:           model.Name,
		description:    model.Description.String,
		tags:           valueobject.NewTags(stringToStringArr(model.Tags)...),
		createdAt:      model.CreatedAt,
		updatedAt:      model.UpdatedAt,
		lastRevisedAt:  model.LastRevisedAt,
		nextRevisionAt: model.NextRevisionAt,
		deletedAt:      pointers.New(model.DeletedAt.Time),
	}, nil
}
