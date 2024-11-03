package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type ChangeName struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

type ChangeNameHandler struct {
	repo reviseitem.Repository
}

func NewChangeNameHandler(repo reviseitem.Repository) *ChangeNameHandler {
	return &ChangeNameHandler{repo: repo}
}

func (h *ChangeNameHandler) Handle(ctx context.Context, cmd ChangeName) error {
	op := errs.Op("application.reviseitem.command.change_name")
	err := h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, errs.
				NewForbiddenError(op, nil, "user is not allowed to modify the item").
				WithMessages([]errs.Message{{Key: "message", Value: "user is not allowed to modify the item"}}).
				WithContext("cmd", cmd)
		}

		err := item.UpdateName(cmd.Name)
		if err != nil {
			return nil, errs.WithOp(op, err, "failed to change name of revise item")
		}

		return item, nil
	})
	if err != nil {
		return errs.WithOp(op, err, "failed to update revise item")
	}
	return nil
}
