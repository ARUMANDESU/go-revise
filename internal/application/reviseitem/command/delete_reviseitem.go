package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type DeleteReviseItem struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type DeleteReviseItemHandler struct {
	repo reviseitem.Repository
}

func NewDeleteReviseItemHandler(repo reviseitem.Repository) *DeleteReviseItemHandler {
	return &DeleteReviseItemHandler{repo: repo}
}

func (h *DeleteReviseItemHandler) Handle(ctx context.Context, cmd DeleteReviseItem) error {
	op := errs.Op("application.reviseitem.command.delete_reviseitem")
	err := h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, errs.
				NewForbiddenError(op, nil, "user is not allowed to modify the item").
				WithMessages([]errs.Message{{Key: "message", Value: "user is not allowed to modify the item"}}).
				WithContext("cmd", cmd)
		}

		item.MarkAsDeleted()

		return item, nil
	})
	if err != nil {
		return errs.WithOp(op, err, "failed to update revise item")
	}

	return nil
}
