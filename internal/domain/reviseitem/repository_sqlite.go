package reviseitem

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
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
	return qtx.UpdateReviseItem(ctx, sqlc.UpdateReviseItemParams{
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

func (r *SqliteRepo) ListUserReviseItems(
	ctx context.Context,
	userID uuid.UUID,
	pagination valueobject.Pagination,
) ([]query.ReviseItem, valueobject.PaginationMetadata, error) {
	q := sqlc.New(r.db)

	reviseItems, err := q.ListUserReviseItems(ctx, sqlc.ListUserReviseItemsParams{
		UserID: userID.String(),
		Limit:  int64(pagination.Limit()),
		Offset: int64(pagination.Offset()),
	})
	if err != nil {
		return nil, valueobject.PaginationMetadata{}, err
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

		// get revisionModels
		revisionModels, err := q.GetRevisionItemRevisions(ctx, item.ID)
		if err != nil {
			return nil, valueobject.PaginationMetadata{}, err
		}
		revisions := make([]time.Time, 0, len(revisionModels))
		for _, revision := range revisionModels {
			revisions = append(revisions, revision.RevisedAt)
		}

		reviseItem.Revisions = revisions

		items = append(items, reviseItem)
	}

	return items, pagination.Metadata(totalCount), nil
}

// --- Query read models implementation ---

func (r *SqliteRepo) GetReviseItem(ctx context.Context, id, userID uuid.UUID) (query.ReviseItem, error) {
	q := sqlc.New(r.db)

	reviseItemModel, err := q.GetReviseItem(ctx, id.String())
	if err != nil {
		return query.ReviseItem{}, err
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
		return query.ReviseItem{}, err
	}
	revisions := make([]time.Time, 0, len(revisionModels))
	for _, revision := range revisionModels {
		revisions = append(revisions, revision.RevisedAt)
	}

	reviseItem.Revisions = revisions

	return reviseItem, nil
}

func (r *SqliteRepo) FetchReviseItemsDueForUser(ctx context.Context, userID uuid.UUID) ([]ReviseItem, error) {
	q := sqlc.New(r.db)

	dayStart := time.Now().Truncate(24 * time.Hour)
	reviseItems, err := q.GetUserReviseItemsByTime(ctx, sqlc.GetUserReviseItemsByTimeParams{
		UserID:         userID.String(),
		NextRevisionAt: dayStart.Add(24 * time.Hour), // before end of day
	})
	if err != nil {
		return nil, err
	}

	var aggregates []ReviseItem
	for _, item := range reviseItems {
		aggregate, err := modelToReviseItem(item)
		if err != nil {
			return nil, err
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
