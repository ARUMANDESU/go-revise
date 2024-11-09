package handler

import (
	"github.com/gofrs/uuid"
	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) RegisterUser(c tb.Context) error {
	op := errs.Op("tgbot.hander.register_user")
	return errs.NewAlreadyExistsError(op, nil, "failed to register new user").
		WithMessages([]errs.Message{{Key: "message", Value: "user already exists"}}).
		WithContext("id", uuid.Must(uuid.NewV7()))
}
