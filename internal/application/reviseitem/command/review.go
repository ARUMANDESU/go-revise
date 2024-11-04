package command

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type Review struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type ReviewHandler struct {
	repo reviseitem.Repository
}

func NewReviewHandler(repo reviseitem.Repository) ReviewHandler {
	return ReviewHandler{repo: repo}
}

func (h *ReviewHandler) Handle(ctx context.Context, cmd Review) error {
	op := errs.Op("application.reviseitem.command.review")
	if cmd.ID.IsNil() {
		return errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "id must be provided").
			WithMessages([]errs.Message{{Key: "message", Value: "id must be provided"}}).
			WithContext("cmd", cmd)
	}
	if cmd.UserID.IsNil() {
		return errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "user_id must be provided").
			WithMessages([]errs.Message{{Key: "message", Value: "user_id must be provided"}}).
			WithContext("cmd", cmd)
	}

	err := h.repo.Update(
		ctx,
		cmd.ID,
		func(ri *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
			if ri.CanModify(cmd.UserID) {
				return nil, errs.
					NewForbiddenError(op, nil, "user is not allowed to modify the item").
					WithMessages([]errs.Message{{Key: "message", Value: "user is not allowed to modify the item"}}).
					WithContext("cmd", cmd)
			}
			ri.Review()
			return ri, nil
		},
	)
	if err != nil {
		return errs.WithOp(op, err, "failed to update revise item")
	}
	return nil
}
