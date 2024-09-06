package tgbot

import (
	"fmt"
	tb "gopkg.in/telebot.v3"
	"strings"
	"unicode/utf8"

	"github.com/ARUMANDESU/go-revise/internal/domain"
)

func DisplayReviewItemMarkdown(item domain.ReviseItem) string {
	return fmt.Sprintf("*%s*\n\n*Tags*: %s\n*Iteration*: %d\n*Next Revision At*: %s",
		item.Name, strings.Join(item.Tags, ", "), item.Iteration, item.NextRevisionAt.Format("2006-01-02 15:04"))
}

func DisplayReviewItemsMarkdown(items []domain.ReviseItem, offset int) string {
	var result strings.Builder

	for i, item := range items {
		result.WriteString(fmt.Sprintf("%d: %s\n", i+1+offset, item.Name))
		result.WriteString(fmt.Sprintf("\t- *Tags*: %s\n", strings.Join(item.Tags, ", ")))
		result.WriteString(fmt.Sprintf("\t- *Iteration*: %d\n", item.Iteration))
		result.WriteString(fmt.Sprintf("\t- *Next Revision At*: %s\n", item.NextRevisionAt.Format("2006-01-02 15:04")))
		result.WriteString("\n")
	}

	return result.String()
}

// ProvideItemButtons provides a list of buttons for each item.
// The offset is used to determine the index of the first item in the list.
//
//	Max 5 items per row.
//	[1] [2] [3] [4] [5]
//	[6] [7] [8] [9] [10]
func (b *Bot) provideItemButtons(ctx tb.Context, items []domain.ReviseItem, offset int) [][]tb.InlineButton {
	buttons := make([][]tb.InlineButton, 0, 1)

	// Add a new row of buttons every 5 items
	for i := 0; i < len(items); i += 5 {
		row := make([]tb.InlineButton, 0, 5)
		for j := 0; j < 5 && i+j < len(items); j++ {
			item := items[i+j]
			row = append(row, tb.InlineButton{
				Text:   fmt.Sprintf("%d", i+j+1+offset),
				Unique: item.ID.String(),
			})
		}
		buttons = append(buttons, row)
	}

	return buttons
}

// SetItemButtonsHandler sets the handler for each item button.
// Buttons are grouped into rows of 5.
func (b *Bot) setItemButtonsHandler(ctx tb.Context, items []domain.ReviseItem, itemButtons [][]tb.InlineButton) {
	itemsMenuButtons := [][]tb.InlineButton{
		{
			BackButton,
		},
	}

	// set handlers for each item button
	for i, row := range itemButtons {
		for j, button := range row {
			item := items[i*5+j]
			b.bot.Handle(&button, b.ReviseItem(item.ID.String(), itemsMenuButtons))
		}
	}

}

func truncateText(text string, maxLength int) string {
	if utf8.RuneCountInString(text) <= maxLength {
		return text
	}
	truncated := text
	runeCount := 0
	for i := 0; i < len(text); i++ {
		if runeCount >= maxLength-3 {
			return string([]rune(truncated)[:maxLength-3]) + "..."
		}
		_, size := utf8.DecodeRuneInString(text[i:])
		i += size - 1
		runeCount++
	}
	return truncated
}
