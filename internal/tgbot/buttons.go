package tgbot

import tb "gopkg.in/telebot.v3"

var (
	EmptyButtonInline = tb.InlineButton{Text: " ", Data: ""}

	// StartButton is a button to start the bot
	StartButtonInline = tb.InlineButton{Text: "/start", Data: "start"}
	// HelpButton is a button to show help message
	HelpButtonInline = tb.InlineButton{Text: "/help", Data: "help"}
	// Revise menu buttons
	ReviseMenuButtonInline   = tb.InlineButton{Text: "Revise Menu", Unique: "revise_menu"}
	ReviseListButtonInline   = tb.InlineButton{Text: "List Review Items", Unique: "revise_list", Data: "revise_list"}
	ReviseCreateButtonInline = tb.InlineButton{Text: "Create Review Item", Unique: "revise_create", Data: "revise_create"}

	// Revise list buttons
	NextButton = tb.InlineButton{Text: "next", Unique: "next"}
	PrevButton = tb.InlineButton{Text: "prev", Unique: "prev"}

	// Revise item buttons
	ResetButton  = tb.InlineButton{Text: "reset", Unique: "reset"}
	BackButton   = tb.InlineButton{Text: "back", Unique: "back"}
	DeleteButton = tb.InlineButton{Text: "delete", Unique: "delete"}
	EditButton   = tb.InlineButton{Text: "edit", Unique: "edit"}
)

var (
	menuButtons = []tb.InlineButton{
		ReviseListButtonInline,
		ReviseCreateButtonInline,
	}
)
