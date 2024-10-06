package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
)

type NewReviseItem struct {
	ID          uuid.UUID        `json:"id"`
	UserID      uuid.UUID        `json:"user_id"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Tags        valueobject.Tags `json:"tags,omitempty"`
}

func (n NewReviseItem) toArgs() reviseitem.NewReviseItemArgs {
	return reviseitem.NewReviseItemArgs{
		ID:          n.ID,
		UserID:      n.UserID,
		Name:        n.Name,
		Description: n.Description,
		Tags:        n.Tags,
	}
}

type NewReviseItemHandler struct {
	repo reviseitem.Repository
}

func NewNewReviseItemHandler(repo reviseitem.Repository) *NewReviseItemHandler {
	return &NewReviseItemHandler{repo: repo}
}

func (h *NewReviseItemHandler) Handle(ctx context.Context, cmd NewReviseItem) error {
	item, err := reviseitem.NewReviseItem(cmd.toArgs())
	if err != nil {
		return err
	}

	if err := h.repo.Save(ctx, *item); err != nil {
		return err
	}

	return nil
}
