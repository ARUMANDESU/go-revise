package handler

import (
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v4"
)

func (h *Handler) StartBot(c tb.Context) error {
	startMsg := strings.Builder{}
	startMsg.WriteString(fmt.Sprintf("Hello, *%s*\\!\n\n", c.Chat().FirstName))
	startMsg.WriteString("👋 *Welcome to Go\\-Revise\\!*\n\n")

	startMsg.WriteString("*What this bot does:*\n")
	startMsg.WriteString("• Help you retain information\n")
	startMsg.WriteString("• Reinforce learning over time\n")
	startMsg.WriteString("• Send revision reminders\n\n")

	startMsg.WriteString("*Getting Started:*\n")
	startMsg.WriteString("1\\. Use /register to create account\n")
	startMsg.WriteString("2\\. Add your study topics\n")
	startMsg.WriteString("3\\. Let the bot handle your revision schedule\n\n")

	startMsg.WriteString("_Ready to enhance your learning journey? Use /register to begin\\!_")

	return c.Send(startMsg.String(), &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
}
