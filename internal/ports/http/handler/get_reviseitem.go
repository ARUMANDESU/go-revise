package handler

import (
	"net/http"

	"github.com/gofrs/uuid"

	reviseitemquery "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httperr"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httpio"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) GetReviseItem(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.get_revise_item")

	var input struct {
		ID     uuid.UUID `json:"id"`
		UserID uuid.UUID `json:"user_id"`
	}
	if err := httpio.ReadJSON(w, r, &input); err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
		return
	}

	query := reviseitemquery.GetReviseItem{
		ID:     input.ID,
		UserID: input.UserID,
	}

	reviseItem, err := h.app.ReviseItem.Query.GetReviseItem.Handle(r.Context(), query)
	if err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to get revise item"))
		return
	}

	httpio.Success(w, r, http.StatusOK, httpio.Envelope{"revise_item": reviseItem})
}
