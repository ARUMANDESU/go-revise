package handler

import (
	"net/http"

	"github.com/gofrs/uuid"

	reviseitemquery "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) GetReviseItem(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.GetReviseItem")

	var input struct {
		ID     uuid.UUID `json:"id"`
		UserID uuid.UUID `json:"user_id"`
	}
	if err := readJSON(w, r, &input); err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
		return
	}

	query := reviseitemquery.GetReviseItem{
		ID:     input.ID,
		UserID: input.UserID,
	}

	reviseItem, err := h.app.ReviseItem.Query.GetReviseItem.Handle(r.Context(), query)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to get revise item"))
		return
	}

	response := map[string]interface{}{"revise_item": reviseItem}

	err = writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to write response"))
		return
	}
}
