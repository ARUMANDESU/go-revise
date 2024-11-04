package handler

import (
	"net/http"

	userquery "github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.GetUser")

	var query userquery.GetUser
	if err := readJSON(w, r, &query); err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
		return
	}

	user, err := h.app.User.Queries.GetUser.Handle(r.Context(), query)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to get user"))
		return
	}

	response := map[string]any{"user": user}

	err = writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to write response"))
		return
	}
}
