package tgbot

import (
	"context"
	"errors"
	"fmt"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	tb "gopkg.in/telebot.v3"
)

func (b *Bot) handleStartCommand(ctx tb.Context) error {
	const op = "tgbot.Bot.handler.handleStartCommand"
	log := b.log.With("op", op)

	_, err := b.UserService.Create(context.Background(), ctx.Sender().ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyExists):
			ctx.Send("Hello again!")
			return nil
		default:
			log.Error("failed to create user", "error", err)
			return err
		}
	}

	ctx.Send(fmt.Sprintf("Hello, %s! Welcome to my bot.", ctx.Sender().FirstName))
	return nil
}

func (b *Bot) handleHelpCommand(ctx tb.Context) error {
	ctx.Send("This is a help message.", &tb.ReplyMarkup{
		ResizeKeyboard: true,
		InlineKeyboard: [][]tb.InlineButton{
			{
				ReviseMenuButtonInline,
			},
		},
	},
	)

	return nil
}

func (b *Bot) handleReviseMenuCommand(ctx tb.Context) error {
	const op = "tgbot.Bot.handler.handleReviseMenuCommand"
	log := b.log.With("op", op)

	log.Debug("Revise menu command received")
	ctx.Edit("Revise commands:", &tb.ReplyMarkup{
		ResizeKeyboard: true,
		InlineKeyboard: [][]tb.InlineButton{
			{
				ReviseListButtonInline,
				ReviseCreateButtonInline,
			},
		},
	})

	return nil
}

func (b *Bot) handleReviseListCommand(ctx tb.Context) error {
	const op = "tgbot.Bot.handler.handleReviseListCommand"
	log := b.log.With("op", op)

	log.Debug("Revise list command received")

	reviseList, _, err := b.ReviseService.List(
		context.Background(),
		domain.ListReviseItemDTO{
			UserID:     ctx.Sender().ID,
			Pagination: domain.NewPagination(1, 10),
			Sort:       domain.DefaultSort(),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			ctx.EditOrSend("You don't have any revise items.", ctx.Message().ReplyMarkup)
			return nil
		default:
			log.Error("failed to list revise items", "error", err)
			ctx.EditOrSend("Failed to list revise items.", ctx.Message().ReplyMarkup)
			return err
		}
	}

	if len(reviseList) == 0 {
		ctx.EditOrSend("You don't have any revise items.", ctx.Message().ReplyMarkup)
		return nil
	}

	log.Debug("Revise list items", "items", reviseList)
	var message string
	for _, item := range reviseList {
		message += fmt.Sprintf("%s\n Last revised: %s\n Next revise: %s\n\n", item.Name, item.LastRevisedAt.Format("2006-01-02"), item.NextRevisionAt.Format("2006-01-02"))
	}

	ctx.Send(message, &tb.ReplyMarkup{
		ResizeKeyboard: true,
		InlineKeyboard: [][]tb.InlineButton{
			{
				tb.InlineButton{Text: "Back", Data: "revise_menu"},
				ReviseCreateButtonInline,
			},
		},
	})

	return nil
}

func (b *Bot) handleReviseCreateCommand(ctx tb.Context) error {
	const op = "tgbot.Bot.handler.handleReviseCreateCommand"
	log := b.log.With("op", op)

	log.Debug("Revise create command received")

	ctx.Send("Enter the title of the item you want to revise.", &tb.ReplyMarkup{
		ResizeKeyboard: true,
		ForceReply:     true,
	})

	ctx.Respond(&tb.CallbackResponse{
		Text: "Enter the title of the item you want to revise.",
	})

	var (
		name        string
		description string
	)

	ctx.Bot().Handle(tb.OnText, func(ctx tb.Context) error {
		name = ctx.Text()

		ctx.Send("Enter the description of the item you want to revise.", &tb.ReplyMarkup{
			ResizeKeyboard: true,
			ForceReply:     true,
		})

		ctx.Respond(&tb.CallbackResponse{
			Text: "Enter the description of the item you want to revise.",
		})

		ctx.Bot().Handle(tb.OnText, func(ctx tb.Context) error {
			description = ctx.Text()

			user, err := b.UserService.GetByChatID(context.Background(), ctx.Chat().ID)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrNotFound):
					ctx.Send("You are not registered. Please use /start command to register.")
					return nil
				default:
					log.Error("failed to get user by chat ID", "error", err)
					return err
				}
			}

			reviseItem, err := b.ReviseService.Create(
				context.Background(),
				domain.CreateReviseItemDTO{
					UserID:      user.ID.String(),
					Name:        name,
					Description: description,
				},
			)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrInvalidArgument):
					ctx.Send(fmt.Sprintf("Failed to create revise item. %s", err))
					return nil
				default:
					log.Error("failed to create revise item", "error", err)
					return err
				}
			}

			ctx.Send(fmt.Sprintf("Revise item created: \nTitle: %s\nDescription: %s", reviseItem.Name, reviseItem.Description))

			return nil
		})

		return nil
	})

	return nil
}
