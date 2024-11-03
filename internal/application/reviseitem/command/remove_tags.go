package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type RemoveTags struct {
	ID     uuid.UUID        `json:"id"`
	UserID uuid.UUID        `json:"user_id"`
	Tags   valueobject.Tags `json:"tags"`
}

type RemoveTagsHandler struct {
	repo reviseitem.Repository
}

func NewRemoveTagsHandler(repo reviseitem.Repository) *RemoveTagsHandler {
	return &RemoveTagsHandler{repo: repo}
}

func (h *RemoveTagsHandler) Handle(ctx context.Context, cmd RemoveTags) error {
	op := errs.Op("application.reviseitem.command.remove_tags")
	err := h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, errs.
				NewForbiddenError(op, nil, "user is not allowed to modify the item").
				WithMessages([]errs.Message{{Key: "message", Value: "user is not allowed to modify the item"}}).
				WithContext("cmd", cmd)
		}

		err := item.RemoveTags(cmd.Tags)
		if err != nil {
			return nil, errs.WithOp(op, err, "failed to remove tags from revise item")
		}

		return item, nil
	})
	if err != nil {
		return errs.WithOp(op, err, "failed to update revise item")
	}
	return nil
}
