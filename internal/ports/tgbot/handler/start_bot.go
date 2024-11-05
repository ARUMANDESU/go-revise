package handler

import (
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v4"
)

func (h *Handler) StartBot(c tb.Context) error {
	startMsg := strings.Builder{}
	startMsg.WriteString(fmt.Sprintf("Hello, %s", c.Chat().FirstName))

	startMsg.WriteString("👋 Welcome to Go-Revise!\n\n")
	startMsg.WriteString("I'm here to help you retain and reinforce information over time. ")
	startMsg.WriteString(
		"Whether it's an article you've read about maps in Go or any other topic, ",
	)
	startMsg.WriteString(
		"I ensure you don’t forget by prompting you to revisit the content at regular intervals.\n\n",
	)
	startMsg.WriteString("To get started, please register by sending the /register command.\n")
	startMsg.WriteString(
		"Once registered, you can start adding topics and keep your knowledge fresh!",
	)
	return c.Send(startMsg.String())
}
