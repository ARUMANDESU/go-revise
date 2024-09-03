package tgbot

import tb "gopkg.in/telebot.v3"

var (
	// StartButton is a button to start the bot
	StartButtonInline = tb.InlineButton{Text: "/start", Data: "start"}
	// HelpButton is a button to show help message
	HelpButtonInline = tb.InlineButton{Text: "/help", Data: "help"}
	// Revise menu buttons
	ReviseMenuButtonInline   = tb.InlineButton{Text: "/revise_menu", Unique: "revise_menu"}
	ReviseListButtonInline   = tb.InlineButton{Text: "/revise_list", Unique: "revise_list"}
	ReviseCreateButtonInline = tb.InlineButton{Text: "/revise_create", Unique: "revise_create"}
)
