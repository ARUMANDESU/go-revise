package command

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
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
	return h.repo.Update(ctx, cmd.ID, func(item *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
		if !item.CanModify(cmd.UserID) {
			return nil, fmt.Errorf("revise item cannot be modified")
		}

		err := item.UpdateName(cmd.Name)
		if err != nil {
			return nil, err
		}

		return item, nil
	})
}
