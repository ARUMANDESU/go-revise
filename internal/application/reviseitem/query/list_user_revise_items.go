package query

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type ListUserReviseItemsReadModel interface {
	ListUserReviseItems(
		ctx context.Context,
		userID uuid.UUID,
		pagination valueobject.Pagination,
	) ([]ReviseItem, valueobject.PaginationMetadata, error)
}

type ListUserReviseItems struct {
	UserID     uuid.UUID  `json:"user_id"`
	Pagination Pagination `json:"pagination"`
}

type ListUserReviseItemsHandler struct {
	readModel ListUserReviseItemsReadModel
}

func NewListUserReviseItemsHandler(
	readModel ListUserReviseItemsReadModel,
) ListUserReviseItemsHandler {
	return ListUserReviseItemsHandler{readModel: readModel}
}

func (h ListUserReviseItemsHandler) Handle(
	ctx context.Context,
	query ListUserReviseItems,
) ([]ReviseItem, valueobject.PaginationMetadata, error) {
	const op = "reviseitem.query.list_user_revise_items"
	if query.UserID.IsNil() {
		return nil, valueobject.PaginationMetadata{}, errs.NewIncorrectInputError(
			op,
			errors.New("user_id must not be nil"),
			"user_id-must-not-be-nil",
		)
	}
	pagination := valueobject.NewPagination(query.Pagination.Page, query.Pagination.PageSize)
	return h.readModel.ListUserReviseItems(ctx, query.UserID, pagination)
}
