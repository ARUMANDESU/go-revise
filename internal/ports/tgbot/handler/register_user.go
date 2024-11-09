package handler

import (
	"context"
	"strings"

	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/internal/application/user/command"
	"github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/internal/ports/tgbot/button"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) RegisterUser(c tb.Context) error {
	confirmMsg := strings.Builder{}
	confirmMsg.WriteString("🔐 *Registration Confirmation*\n\n")
	confirmMsg.WriteString("*Data We Store:*\n")
	confirmMsg.WriteString("• Your Telegram Chat ID\n")
	confirmMsg.WriteString("• Revision items you create:\n")
	confirmMsg.WriteString("  \\- Item names\n")
	confirmMsg.WriteString("  \\- Descriptions\n")
	confirmMsg.WriteString("  \\- Custom tags\n")
	confirmMsg.WriteString("  \\- Creation dates\n")
	confirmMsg.WriteString("  \\- Revision schedules\n\n")

	return c.Send(
		confirmMsg.String(),
		&tb.SendOptions{ParseMode: tb.ModeMarkdownV2},
		&tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{button.RegistrationConfirmI}}},
	)
}

func (h *Handler) RegisterUserConfirmed(c tb.Context) error {
	op := errs.Op("tgbot.handler.register_user_confirmed")

	err := h.app.User.Commands.RegisterUser.Handle(
		context.TODO(),
		command.RegisterUser{ChatID: user.TelegramID(c.Chat().ID)},
	)
	if err != nil {
		return errs.WithOp(op, err, "failed to register user")
	}

	err = c.Edit(
		"✅ *Registration Confirmed*",
		&tb.SendOptions{ParseMode: tb.ModeMarkdownV2},
	)
	if err != nil {
		return errs.WithOp(op, err, "failed to update confirmation message")
	}

	msg := strings.Builder{}
	msg.WriteString("🎉 *Welcome\\!* Your registration is complete\\.\n\n")
	msg.WriteString("*What you can do now:*\n")
	msg.WriteString("• Use /help to see available commands\n")
	msg.WriteString("• Set up your preferences with /settings\n")
	msg.WriteString("• Start exploring with /menu\n\n")

	return c.Send(msg.String(), &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
}
