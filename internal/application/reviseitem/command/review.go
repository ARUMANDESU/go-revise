package command

import (
	"context"
	"errors"

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

func NewReviewHandler(repo reviseitem.Repository) *ReviewHandler {
	return &ReviewHandler{repo: repo}
}

func (h *ReviewHandler) Handle(ctx context.Context, cmd Review) error {
	const op = "reviseitem.command.review"
	if cmd.ID.IsNil() {
		return errs.NewIncorrectInputError(
			op,
			errors.New("revise item id must be provided"),
			"revise-item-id-must-be-provided",
		)
	}
	if cmd.UserID.IsNil() {
		return errs.NewIncorrectInputError(
			op,
			errors.New("user_id must be provided"),
			"user_id-must-be-provided",
		)
	}

	return h.repo.Update(
		ctx,
		cmd.ID,
		func(ri *reviseitem.Aggregate) (*reviseitem.Aggregate, error) {
			if ri.CanModify(cmd.UserID) {
				return nil, errs.NewAuthorizationError(
					op,
					errors.New("not authorized to review"),
					"not-authorized-to-review",
				)
			}
			ri.Review()
			return ri, nil
		},
	)
}
