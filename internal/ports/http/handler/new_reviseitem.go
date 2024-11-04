package handler

import (
	"net/http"

	"github.com/gofrs/uuid"

	reviseitemcmd "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/command"
	reviseitemquery "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) NewReviseItem(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.NewReviseItem")
	var input struct {
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name"`
		Description *string   `json:"description,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
	}

	if err := readJSON(w, r, &input); err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
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
		handleError(w, r, errs.WithOp(op, err, "failed to create new revise item"))
		return
	}

	queryReviseItem, err := h.app.ReviseItem.Query.GetReviseItem.Handle(
		r.Context(),
		reviseitemquery.GetReviseItem{ID: reviseItemID, UserID: input.UserID},
	)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to get revise item"))
		return
	}

	response := map[string]interface{}{"revise_item": queryReviseItem}

	err = writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to write response"))
		return
	}
}
