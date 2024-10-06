package command

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
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
	return h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, fmt.Errorf("revise item cannot be modified")
		}

		err := item.RemoveTags(cmd.Tags)
		if err != nil {
			return nil, err
		}

		return item, nil
	})
}
