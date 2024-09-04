package tgbot

import tb "gopkg.in/telebot.v3"

var (
	// StartButton is a button to start the bot
	StartButtonInline = tb.InlineButton{Text: "/start", Data: "start"}
	// HelpButton is a button to show help message
	HelpButtonInline = tb.InlineButton{Text: "/help", Data: "help"}
	// Revise menu buttons
	ReviseMenuButtonInline   = tb.InlineButton{Text: "revise menu", Unique: "revise_menu"}
	ReviseListButtonInline   = tb.InlineButton{Text: "revise list", Unique: "revise_list"}
	ReviseCreateButtonInline = tb.InlineButton{Text: "create new revise item", Unique: "revise_create"}

	// Revise list buttons
	ReviseListNextButtonInline = tb.InlineButton{Text: "next", Unique: "revise_list_next"}
	ReviseListPrevButtonInline = tb.InlineButton{Text: "prev", Unique: "revise_list_prev"}
)
