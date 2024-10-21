package query

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type GetReviseItemReadModel interface {
	GetReviseItem(ctx context.Context, id, userID uuid.UUID) (ReviseItem, error)
}

type GetReviseItem struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type GetReviseItemHandler struct {
	readModel GetReviseItemReadModel
}

func NewGetReviseItemHandler(readModel GetReviseItemReadModel) GetReviseItemHandler {
	return GetReviseItemHandler{
		readModel: readModel,
	}
}

func (h GetReviseItemHandler) Handle(ctx context.Context, query GetReviseItem) (ReviseItem, error) {
	if query.UserID.IsNil() {
		return ReviseItem{}, errs.NewIncorrectInputError(
			"user_id must not be nil",
			"user_id-must-not-be-nil",
		)
	}
	if query.ID.IsNil() {
		return ReviseItem{}, errs.NewIncorrectInputError("id must not be nil", "id-must-not-be-nil")
	}
	return h.readModel.GetReviseItem(ctx, query.ID, query.UserID)
}
