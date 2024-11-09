package handler

import (
	"context"
	"strings"
	"unicode"

	"github.com/gofrs/uuid"
	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/internal/application/reviseitem/command"
	reviseitemquery "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
	"github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (h *Handler) CreateItem(c tb.Context) error {
	op := errs.Op("handler.create_item")

	fullText := c.Message().Text
	commandEnd := strings.Index(fullText, "/revise_create") + len("/revise_create")
	if commandEnd >= len(fullText) {
		return c.Reply(
			"⚠️ *Usage:*\n"+
				"/revise\\_create \"name\" \\[\"description\"\\] \\[\"tags\"\\]\n\n"+
				"*Example:*\n"+
				"/revise\\_create \"Go Maps\" \"Understanding Go maps\" \"go,data\\-structures\"\n\n"+
				"Note:\n"+
				"• *Name* is required\n"+
				"• *Description* and *tags* are optional\n"+
				"• Use quotes for values with spaces\n"+
				"• Tags should be comma\\-separated",
			&tb.SendOptions{ParseMode: tb.ModeMarkdownV2},
		)
	}

	// Parse quoted arguments
	args := parseQuotedArgs(fullText[commandEnd:])
	if len(args) == 0 {
		return c.Reply(
			"❌ *Name is required*",
			&tb.SendOptions{ParseMode: tb.ModeMarkdownV2},
		)
	}

	name := args[0]
	var description string
	var tags []string

	if len(args) > 1 {
		description = args[1]
	}

	if len(args) > 2 {
		tagList := strings.Split(args[2], ",")
		tags = make([]string, 0, len(tagList))
		for _, tag := range tagList {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	queryUser, err := h.app.User.Queries.GetUser.Handle(
		context.TODO(),
		query.GetUser{ChatID: user.TelegramID(c.Chat().ID)},
	)
	if err != nil {
		return errs.WithOp(op, err, "failed to get user")
	}

	userID, err := uuid.FromString(queryUser.ID)
	if err != nil {
		return errs.WithOp(op, err, "failed to parse user ID")
	}

	err = h.app.ReviseItem.Command.NewReviseItem.Handle(
		context.TODO(),
		command.NewReviseItem{
			ID:          reviseitem.NewReviseItemID(),
			UserID:      userID,
			Name:        name,
			Description: description,
			Tags:        valueobject.NewTags(tags...),
		},
	)
	if err != nil {
		return errs.WithOp(op, err, "failed to create item")
	}

	revisionItem, err := h.app.ReviseItem.Query.GetReviseItem.Handle(
		context.TODO(),
		reviseitemquery.GetReviseItem{ID: reviseitem.NewReviseItemID(), UserID: userID},
	)
	if err != nil {
		return errs.WithOp(op, err, "failed to get revision item")
	}

	msg := strings.Builder{}
	msg.WriteString("✅ *Revision Item Created*\n\n")
	msg.WriteString("*Name:* " + escapeMarkdown(revisionItem.Name) + "\n")
	if description != "" {
		msg.WriteString("*Description:* " + escapeMarkdown(revisionItem.Description) + "\n")
	}
	if len(tags) > 0 {
		msg.WriteString(
			"*Tags:* " + escapeMarkdown(strings.Join(revisionItem.Tags.StringArray(), ", ")) + "\n",
		)
	}
	msg.WriteString("\nUse /list to see all your items")

	return c.Reply(msg.String(), &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
}

func escapeMarkdown(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(",
		")", "~", "`", ">", "#",
		"+", "-", "=", "|", "{",
		"}", ".", "!",
	}
	escaped := text
	for _, char := range specialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

func parseQuotedArgs(s string) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false

	// Trim leading spaces
	s = strings.TrimSpace(s)

	for i := 0; i < len(s); i++ {
		char := rune(s[i])

		switch {
		case char == '"':
			// Toggle quote state
			inQuotes = !inQuotes

			// If we're ending quotes, add the argument
			if !inQuotes && currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}

		case inQuotes:
			// If we're inside quotes, add character to current argument
			currentArg.WriteRune(char)

		case !unicode.IsSpace(char):
			// If we're outside quotes and it's not a space,
			// treat it as part of an unquoted argument
			currentArg.WriteRune(char)

			// If it's the last character or next is space, end the argument
			if i == len(s)-1 || unicode.IsSpace(rune(s[i+1])) {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		}
	}

	// Add final argument if any
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}
