package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type AddTags struct {
	ID     uuid.UUID        `json:"id"`
	UserID uuid.UUID        `json:"user_id"`
	Tags   valueobject.Tags `json:"tags"`
}

type AddTagsHandler struct {
	repo reviseitem.Repository
}

func NewAddTagsHandler(repo reviseitem.Repository) *AddTagsHandler {
	return &AddTagsHandler{repo: repo}
}

func (h *AddTagsHandler) Handle(ctx context.Context, cmd AddTags) error {
	op := errs.Op("application.reviseitem.command.add_tags")
	err := h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, errs.
				NewForbiddenError(op, nil, "user is not allowed to modify the item").
				WithMessages([]errs.Message{{Key: "message", Value: "user is not allowed to modify the item"}}).
				WithContext("cmd", cmd)
		}

		err := item.AddTags(cmd.Tags)
		if err != nil {
			return nil, errs.WithOp(op, err, "failed to add tags to revise item")
		}

		return item, nil
	})
	if err != nil {
		return errs.WithOp(op, err, "failed to update revise item")
	}
	return nil
}
