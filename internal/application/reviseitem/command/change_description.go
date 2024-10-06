package command

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
)

type ChangeDescription struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Description string    `json:"description"`
}

type ChangeDescriptionHandler struct {
	repo reviseitem.Repository
}

func NewChangeDescriptionHandler(repo reviseitem.Repository) *ChangeDescriptionHandler {
	return &ChangeDescriptionHandler{repo: repo}
}

func (h *ChangeDescriptionHandler) Handle(ctx context.Context, cmd ChangeDescription) error {
	return h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, fmt.Errorf("revise item cannot be modified")
		}

		err := item.UpdateDescription(cmd.Description)
		if err != nil {
			return nil, err
		}

		return item, nil
	})
}
