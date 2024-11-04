package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
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

func NewNewReviseItemHandler(repo reviseitem.Repository) NewReviseItemHandler {
	return NewReviseItemHandler{repo: repo}
}

func (h *NewReviseItemHandler) Handle(ctx context.Context, cmd NewReviseItem) error {
	op := errs.Op("application.reviseitem.command.new_reviseitem")
	item, err := reviseitem.NewReviseItem(cmd.toArgs())
	if err != nil {
		return errs.WithOp(op, err, "failed to create new revise item")
	}

	if err := h.repo.Save(ctx, *reviseitem.NewAggregate(item)); err != nil {
		return errs.WithOp(op, err, "failed to save new revise item")
	}

	return nil
}
