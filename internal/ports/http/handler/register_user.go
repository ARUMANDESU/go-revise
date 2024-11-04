package handler

import (
	"net/http"

	usercmd "github.com/ARUMANDESU/go-revise/internal/application/user/command"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	op := errs.Op("handler.RegisterUser")

	var cmd usercmd.RegisterUser
	if err := readJSON(w, r, &cmd); err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to read JSON"))
		return
	}

	err := h.app.User.Commands.RegisterUser.Handle(r.Context(), cmd)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to register user"))
		return
	}

	err = writeJSON(w, http.StatusCreated, nil, nil)
	if err != nil {
		handleError(w, r, errs.WithOp(op, err, "failed to write response"))
		return
	}
}
