package handler

import (
	"context"
	"strings"

	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/internal/application/user/command"
	"github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) RegisterUser(c tb.Context) error {
	op := errs.Op("tgbot.hander.register_user")

	err := h.app.User.Commands.RegisterUser.Handle(
		context.TODO(),
		command.RegisterUser{ChatID: user.TelegramID(c.Chat().ID)},
	)
	if err != nil {
		return errs.WithOp(op, err, "failed to register user")
	}

	msg := strings.Builder{}
	msg.WriteString("🎉 *Welcome\\!* Your registration is complete\\.\n\n")
	msg.WriteString("*What you can do now:*\n")
	msg.WriteString("• Use /help to see available commands\n")
	msg.WriteString("• Set up your preferences with /settings\n")
	msg.WriteString("• Start exploring with /menu\n\n")

	return c.Send(msg.String(), &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
}
