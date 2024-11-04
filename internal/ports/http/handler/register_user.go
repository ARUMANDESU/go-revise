package handler

import (
	"net/http"

	usercmd "github.com/ARUMANDESU/go-revise/internal/application/user/command"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var cmd usercmd.RegisterUser
	if err := readJSON(w, r, &cmd); err != nil {
		handleError(w, r, err)
		return
	}

	err := h.app.User.Commands.RegisterUser.Handle(r.Context(), cmd)
	if err != nil {
		handleError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusCreated, nil, nil)
	if err != nil {
		handleError(w, r, err)
		return
	}
}
