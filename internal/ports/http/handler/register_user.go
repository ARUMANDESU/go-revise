package handler

import (
	"net/http"

	usercmd "github.com/ARUMANDESU/go-revise/internal/application/user/command"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httperr"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/httpio"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.register_user")

	var cmd usercmd.RegisterUser
	if err := httpio.ReadJSON(w, r, &cmd); err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
		return
	}

	err := h.app.User.Commands.RegisterUser.Handle(r.Context(), cmd)
	if err != nil {
		httperr.HandleError(w, r, errs.WithOp(op, err, "failed to register user"))
		return
	}

	httpio.Success(w, r, http.StatusCreated, nil)
}
