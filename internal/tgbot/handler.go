package tgbot

import (
	"fmt"

	tb "gopkg.in/telebot.v3"
)

func (b *Bot) handleStartCommand(ctx tb.Context) error {
	ctx.Send(fmt.Sprintf("Hello, %s! Welcome to my bot.", ctx.Sender().FirstName))
	return nil
}

func (b *Bot) handleHelpCommand(ctx tb.Context) error {
	ctx.Send("This is a help message.", &tb.ReplyMarkup{
		ResizeKeyboard: true,
		InlineKeyboard: [][]tb.InlineButton{
			{
				tb.InlineButton{Text: "/start", InlineQueryChat: "/start"},
				tb.InlineButton{Text: "/help", Unique: "/help"},
				tb.InlineButton{Text: "/revise_list", Unique: "/revise_list"},
				tb.InlineButton{Text: "/revise_create", Unique: "/revise_create"},
			},
		},
	},
	)

	return nil
}
