package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type ChangeDescription struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Description string    `json:"description"`
}

type ChangeDescriptionHandler struct {
	repo reviseitem.Repository
}

func NewChangeDescriptionHandler(repo reviseitem.Repository) ChangeDescriptionHandler {
	return ChangeDescriptionHandler{repo: repo}
}

func (h *ChangeDescriptionHandler) Handle(ctx context.Context, cmd ChangeDescription) error {
	op := errs.Op("application.reviseitem.command.change_description")
	err := h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, errs.
				NewForbiddenError(op, nil, "user is not allowed to modify the item").
				WithMessages([]errs.Message{{Key: "message", Value: "user is not allowed to modify the item"}}).
				WithContext("cmd", cmd)
		}

		err := item.UpdateDescription(cmd.Description)
		if err != nil {
			return nil, errs.WithOp(op, err, "failed to change description of revise item")
		}

		return item, nil
	})
	if err != nil {
		return errs.WithOp(op, err, "failed to update revise item")
	}
	return nil
}
