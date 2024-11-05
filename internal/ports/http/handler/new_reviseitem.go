package handler

import (
	"net/http"

	"github.com/gofrs/uuid"

	reviseitemcmd "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/command"
	reviseitemquery "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httperr"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httpio"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) NewReviseItem(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.new_revise_item")
	var input struct {
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name"`
		Description *string   `json:"description,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
	}

	if err := httpio.ReadJSON(w, r, &input); err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
		return
	}

	reviseItemID := uuid.Must(uuid.NewV4())
	cmd := reviseitemcmd.NewReviseItem{
		ID:     reviseItemID,
		UserID: input.UserID,
		Name:   input.Name,
		Tags:   valueobject.NewTags(input.Tags...),
	}
	if input.Description != nil {
		cmd.Description = *input.Description
	}

	err := h.app.ReviseItem.Command.NewReviseItem.Handle(r.Context(), cmd)
	if err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to create new revise item"))
		return
	}

	queryReviseItem, err := h.app.ReviseItem.Query.GetReviseItem.Handle(
		r.Context(),
		reviseitemquery.GetReviseItem{ID: reviseItemID, UserID: input.UserID},
	)
	if err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to get revise item"))
		return
	}

	httpio.Success(w, r, http.StatusCreated, httpio.Envelope{"revise_item": queryReviseItem})
}
