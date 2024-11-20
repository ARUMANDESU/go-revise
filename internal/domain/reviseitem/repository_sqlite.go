package reviseitem

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/adapters/db/sqlc"
	"github.com/ARUMANDESU/go-revise/internal/adapters/db/sqliterr"
	"github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
	"github.com/ARUMANDESU/go-revise/pkg/pointers"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(db *sql.DB) SQLiteRepo {
	return SQLiteRepo{db: db}
}

// Save saves a revise item.
func (r *SQLiteRepo) Save(ctx context.Context, item Aggregate) (_ error) {
	op := errs.Op("domain.reviseitem.sqlite.save")
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

	err := q.SaveReviseItem(ctx, args)
	if err != nil {
		return sqliterr.Handle(op, err, "failed to save revise item").WithContext("args", args)
	}

	return nil
}

func (r *SQLiteRepo) withTx(ctx context.Context, op errs.Op, fn func(*sqlc.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return sqliterr.HandleTx(op, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.
					With(slog.String("op", string(op))).
					Error("failed to rollback transaction",
						logutil.Err(rollbackErr),
						"original_error", err)
			}
		}
	}()

	qtx := sqlc.New(tx)
	if err = fn(qtx); err != nil {
		return err // Already wrapped with operation
	}

	if err = tx.Commit(); err != nil {
		return sqliterr.HandleTx(op, err, "failed to commit transaction")
	}

	return nil
}

// Update updates a revise item.
func (r *SQLiteRepo) Update(ctx context.Context, id uuid.UUID, fn UpdateFn) (err error) {
	op := errs.Op("domain.reviseitem.sqlite.update")

	return r.withTx(ctx, op, func(q *sqlc.Queries) error {
		reviseItemModel, err := q.GetReviseItem(ctx, id.String())
		if err != nil {
			return sqliterr.Handle(op, err, "failed to get revise item").WithContext("id", id)
		}

		reviseItem, err := modelToReviseItem(reviseItemModel)
		if err != nil {
			return errs.WithOp(op, err, "failed to convert model to revise item")
		}

		aggregate := NewAggregate(&reviseItem)

		aggregate, err = fn(aggregate)
		if err != nil {
			return errs.WithOp(op, err, "failed to update revise item")
		}

		for _, r := range aggregate.Revisions() {
			args := sqlc.CreateRevisionParams{
				ID:           r.ID().String(),
				ReviseItemID: aggregate.ID().String(),
				RevisedAt:    r.RevisedAt(),
			}

			err = q.CreateRevision(ctx, args)
			if err != nil {
				return sqliterr.
					Handle(op, err, "failed to create revision").
					WithContext("args", args)
			}
		}

		tags := aggregate.Tags()
		err = q.UpdateReviseItem(ctx, sqlc.UpdateReviseItemParams{
			Name: aggregate.Name(),
			Description: sql.NullString{
				String: aggregate.Description(),
				Valid:  aggregate.Description() != "",
			},
			Tags:           stringArrToString(tags.StringArray()),
			CreatedAt:      aggregate.CreatedAt(),
			UpdatedAt:      aggregate.UpdatedAt(),
			LastRevisedAt:  aggregate.LastRevisedAt(),
			NextRevisionAt: aggregate.NextRevisionAt(),
			ID:             aggregate.ID().String(),
		})
		if err != nil {
			return sqliterr.
				Handle(op, err, "failed to update revise item").
				WithContext("aggregate", aggregate)
		}

		return nil
	})
}

func (r *SQLiteRepo) ListUserReviseItems(
	ctx context.Context,
	userID uuid.UUID,
	pagination valueobject.Pagination,
) ([]query.ReviseItem, valueobject.PaginationMetadata, error) {
	op := errs.Op("domain.reviseitem.sqlite.list_user_revise_items")
	q := sqlc.New(r.db)

	reviseItems, err := q.ListUserReviseItems(ctx, sqlc.ListUserReviseItemsParams{
		UserID: userID.String(),
		Limit:  int64(pagination.Limit()),
		Offset: int64(pagination.Offset()),
	})
	if err != nil {
		return nil, valueobject.PaginationMetadata{}, sqliterr.Handle(
			op,
			err,
			"failed to list user revise items",
		)
	}

	var (
		items      []query.ReviseItem
		totalCount int
	)
	for _, item := range reviseItems {
		totalCount = int(item.Count)
		var deletedAt *time.Time
		if item.DeletedAt.Valid {
			deletedAt = pointers.New(item.DeletedAt.Time)
		}
		reviseItem := query.ReviseItem{
			ID:             uuid.FromStringOrNil(item.ID),
			UserID:         uuid.FromStringOrNil(item.UserID),
			Name:           item.Name,
			Description:    item.Description.String,
			Tags:           valueobject.NewTags(stringToStringArr(item.Tags)...),
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
			DeletedAt:      deletedAt,
			NextRevisionAt: item.NextRevisionAt,
			LastRevisedAt:  item.LastRevisedAt,
			Revisions:      nil,
		}

		revisions, err := r.getRevisions(ctx, q, item.ID)
		if err != nil {
			return nil, valueobject.PaginationMetadata{}, sqliterr.
				Handle(op, err, "failed to get revisions").
				WithContext("id", item.ID)
		}

		reviseItem.Revisions = revisions

		items = append(items, reviseItem)
	}

	return items, pagination.Metadata(totalCount), nil
}

func (r *SQLiteRepo) getRevisions(
	ctx context.Context,
	q *sqlc.Queries,
	reviseItemID string,
) ([]time.Time, error) {
	op := errs.Op("domain.reviseitem.sqlite.get_revisions")
	revisionModels, err := q.GetRevisionItemRevisions(ctx, reviseItemID)
	if err != nil {
		return nil, sqliterr.
			Handle(op, err, "failed to get revision item revisions").
			WithContext("id", reviseItemID)
	}
	revisions := make([]time.Time, 0, len(revisionModels))
	for _, revision := range revisionModels {
		revisions = append(revisions, revision.RevisedAt)
	}
	return revisions, nil
}

// --- Query read models implementation ---

func (r *SQLiteRepo) GetReviseItem(
	ctx context.Context,
	id, userID uuid.UUID,
) (query.ReviseItem, error) {
	op := errs.Op("domain.reviseitem.sqlite.get_revise_item")
	q := sqlc.New(r.db)

	reviseItemModel, err := q.GetReviseItem(ctx, id.String())
	if err != nil {
		return query.ReviseItem{}, sqliterr.
			Handle(op, err, "failed to get revise item").
			WithContext("id", id)
	}

	var deletedAt *time.Time
	if reviseItemModel.DeletedAt.Valid {
		deletedAt = pointers.New(reviseItemModel.DeletedAt.Time)
	}
	reviseItem := query.ReviseItem{
		ID:             uuid.FromStringOrNil(reviseItemModel.ID),
		UserID:         uuid.FromStringOrNil(reviseItemModel.UserID),
		Name:           reviseItemModel.Name,
		Description:    reviseItemModel.Description.String,
		Tags:           valueobject.NewTags(stringToStringArr(reviseItemModel.Tags)...),
		CreatedAt:      reviseItemModel.CreatedAt,
		UpdatedAt:      reviseItemModel.UpdatedAt,
		DeletedAt:      deletedAt,
		NextRevisionAt: reviseItemModel.NextRevisionAt,
		LastRevisedAt:  reviseItemModel.LastRevisedAt,
		Revisions:      nil,
	}

	// get revisionModels
	revisionModels, err := q.GetRevisionItemRevisions(ctx, id.String())
	if err != nil {
		return query.ReviseItem{}, sqliterr.
			Handle(op, err, "failed to get revise item revisions").
			WithContext("id", id)
	}
	revisions := make([]time.Time, 0, len(revisionModels))
	for _, revision := range revisionModels {
		revisions = append(revisions, revision.RevisedAt)
	}

	reviseItem.Revisions = revisions

	return reviseItem, nil
}

func (r *SQLiteRepo) FetchReviseItemsDueForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]ReviseItem, error) {
	op := errs.Op("domain.reviseitem.sqlite.fetch_revise_items_due_for_user")
	q := sqlc.New(r.db)

	dayStart := time.Now().Truncate(24 * time.Hour)
	reviseItems, err := q.GetUserReviseItemsByTime(ctx, sqlc.GetUserReviseItemsByTimeParams{
		UserID:         userID.String(),
		NextRevisionAt: dayStart.Add(24 * time.Hour), // before end of day
	})
	if err != nil {
		return nil, sqliterr.
			Handle(op, err, "failed to get user revise items by time").
			WithContext("user_id", userID)
	}

	var aggregates []ReviseItem
	for _, item := range reviseItems {
		aggregate, err := modelToReviseItem(item)
		if err != nil {
			return nil, errs.WithOp(op, err, "failed to convert model to revise item")
		}
		aggregates = append(aggregates, aggregate)
	}

	return aggregates, nil
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
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = pointers.New(model.DeletedAt.Time)
	}
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
		deletedAt:      deletedAt,
	}, nil
}
