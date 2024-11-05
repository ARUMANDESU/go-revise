package handler

import (
	"net/http"

	userquery "github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httperr"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httpio"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.get_user")

	var query userquery.GetUser
	if err := httpio.ReadJSON(w, r, &query); err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
		return
	}

	user, err := h.app.User.Queries.GetUser.Handle(r.Context(), query)
	if err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to get user"))
		return
	}

	httpio.Success(w, r, http.StatusOK, httpio.Envelope{"user": user})
}
