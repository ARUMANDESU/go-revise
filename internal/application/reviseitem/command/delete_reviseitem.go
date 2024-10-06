package command

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
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
	err := h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanBeDeleted(cmd.UserID) {
			return nil, fmt.Errorf("revise item cannot be deleted")
		}

		item.MarkAsDeleted()

		return item, nil
	})
	if err != nil {
		return err
	}

	return nil
}
