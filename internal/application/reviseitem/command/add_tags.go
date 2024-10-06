package command

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
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
	return h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, fmt.Errorf("revise item cannot be modified")
		}

		err := item.AddTags(cmd.Tags)
		if err != nil {
			return nil, err
		}

		return item, nil
	})
}
