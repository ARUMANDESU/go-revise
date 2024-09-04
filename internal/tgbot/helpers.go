package tgbot

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ARUMANDESU/go-revise/internal/domain"
)

func DisplayReviewItemsTable(items []domain.ReviseItem) string {
	// Determine column widths
	nameWidth := 15
	tagsWidth := 15
	nextRevWidth := 18

	// Table header
	header := fmt.Sprintf("+%-15s+%-15s+%-18s+\n",
		strings.Repeat("-", nameWidth),
		strings.Repeat("-", tagsWidth),
		strings.Repeat("-", nextRevWidth),
	)
	header += fmt.Sprintf("| %-15s | %-15s | %-218s |\n",
		"Name", "Tags", "Next Revision At")
	header += fmt.Sprintf("+%-15s+%-15s+%-18s+\n",
		strings.Repeat("-", nameWidth),
		strings.Repeat("-", tagsWidth),
		strings.Repeat("-", nextRevWidth),
	)

	// Table rows
	var rows strings.Builder
	for _, item := range items {
		rows.WriteString(fmt.Sprintf("| %-15s | %-15s | %-18s |\n",
			truncateText(item.Name, nameWidth),
			truncateText(fmt.Sprint(item.Tags), tagsWidth),
			item.NextRevisionAt.Format("2006-01-02 15:04"),
		))
	}

	// Bottom border
	footer := fmt.Sprintf("+%-15s+%-15s+%-18s+\n",
		strings.Repeat("-", nameWidth),
		strings.Repeat("-", tagsWidth),
		strings.Repeat("-", nextRevWidth),
	)

	return header + rows.String() + footer
}

func DisplayReviewItemsMarkdown(items []domain.ReviseItem) string {
	var result strings.Builder

	for _, item := range items {
		result.WriteString(fmt.Sprintf("%s\n", item.Name))
		result.WriteString(fmt.Sprintf("\t- *Tags*: %s\n", strings.Join(item.Tags, ", ")))
		result.WriteString(fmt.Sprintf("\t- *Iteration*: %d\n", item.Iteration))
		result.WriteString(fmt.Sprintf("\t- *Next Revision At*: %s\n", item.NextRevisionAt.Format("2006-01-02 15:04")))
		result.WriteString("\n")
	}

	return result.String()
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
