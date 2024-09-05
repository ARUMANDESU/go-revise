package tgbot

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

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
	const op = "tgbot.Bot.handler.handleHelpCommand"
	log := b.log.With("op", op)

	log.Debug("Help command received")

	var helpMessage strings.Builder
	helpMessage.WriteString("*Help menu*:\n")
	helpMessage.WriteString("Here are the available commands:\n")
	helpMessage.WriteString("\n*General commands:*\n")
	helpMessage.WriteString("*/start* - Start the bot\n")
	helpMessage.WriteString("*/help* - Show this help message\n")
	helpMessage.WriteString("\n*Revise commands:*\n")
	helpMessage.WriteString("*/revise_menu* - Revise commands\n")
	helpMessage.WriteString("*/revise_list* - List all revise items\n")
	helpMessage.WriteString("*/revise_create* - Create a new revise item\n")
	helpMessage.WriteString("You can also use the following buttons to navigate:\n")

	ctx.Send(helpMessage.String(),
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		},
		&tb.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]tb.InlineButton{
				{
					ReviseMenuButtonInline,
					ReviseListButtonInline,
					ReviseCreateButtonInline,
				},
			},
		})

	return nil
}

func (b *Bot) handleReviseMenuCommand(ctx tb.Context) error {
	const op = "tgbot.Bot.handler.handleReviseMenuCommand"
	log := b.log.With("op", op)

	log.Debug("Revise menu command received")
	ctx.EditOrSend("Revise commands:", &tb.ReplyMarkup{
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

	var (
		currentPage = 1
		pageSize    = 5
		lastPage    = 1
	)

	paginationButtons := func() []tb.InlineButton {
		var buttons []tb.InlineButton

		if currentPage > 1 {
			buttons = append(buttons, PrevButton)
		} else {
			buttons = append(buttons, EmptyButtonInline)
		}

		buttons = append(buttons, tb.InlineButton{Text: fmt.Sprintf("%d/%d", currentPage, lastPage)})

		if currentPage < lastPage {
			buttons = append(buttons, NextButton)
		} else {
			buttons = append(buttons, EmptyButtonInline)
		}

		return buttons
	}

	displayList := func(pagination *domain.Pagination, sort *domain.Sort) error {
		reviseList, metadata, err := b.ReviseService.List(
			context.Background(),
			domain.ListReviseItemDTO{
				UserID:     ctx.Sender().ID,
				Pagination: pagination,
				Sort:       sort,
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

		currentPage, lastPage = metadata.CurrentPage, metadata.LastPage

		log.Debug("Revise items listed", "count", len(reviseList))
		message := DisplayReviewItemsMarkdown(reviseList)

		log.Debug("Sending revise items list")

		ctx.EditOrSend(message,
			&tb.SendOptions{
				ParseMode: tb.ModeMarkdown,
			},
			&tb.ReplyMarkup{
				ResizeKeyboard: true,
				InlineKeyboard: [][]tb.InlineButton{
					{
						ReviseMenuButtonInline,
						ReviseCreateButtonInline,
					},
					paginationButtons(),
				},
			})

		log.Debug("Revise items list sent")

		return nil
	}

	// initial list
	displayList(domain.NewPagination(1, pageSize), domain.DefaultSort())

	ctx.Bot().Handle(&NextButton, func(ctx tb.Context) error {
		if currentPage == lastPage {
			ctx.Respond(&tb.CallbackResponse{Text: "You are already on the last page."})
			return nil
		}

		displayList(domain.NewPagination(currentPage+1, pageSize), domain.DefaultSort())

		return nil
	})

	ctx.Bot().Handle(&PrevButton, func(ctx tb.Context) error {
		if currentPage == 1 {
			ctx.Respond(&tb.CallbackResponse{Text: "You are already on the first page."})
			return nil
		}

		displayList(domain.NewPagination(currentPage-1, pageSize), domain.DefaultSort())

		return nil
	})

	return nil
}

func (b *Bot) handleReviseCreateCommand(ctx tb.Context) error {
	const op = "tgbot.Bot.handler.handleReviseCreateCommand"
	log := b.log.With("op", op)

	log.Debug("Revise create command received")

	ctx.Send("Enter the title of the item you want to revise.", &tb.ReplyMarkup{
		ResizeKeyboard: true,
	})

	ctx.Respond(&tb.CallbackResponse{
		Text: "Enter the title of the item you want to revise.",
	})

	var (
		wg          sync.WaitGroup
		name        string
		description string
	)

	wg.Add(1)
	ctx.Bot().Handle(tb.OnText, func(ctx tb.Context) error {
		defer wg.Done()
		name = ctx.Text()
		return nil
	})

	wg.Wait()

	ctx.Send("Enter the description of the item you want to revise.", &tb.ReplyMarkup{
		ResizeKeyboard: true,
	})

	ctx.Respond(&tb.CallbackResponse{
		Text: "Enter the description of the item you want to revise.",
	})

	wg.Add(1)
	ctx.Bot().Handle(tb.OnText, func(ctx tb.Context) error {
		defer wg.Done()
		description = ctx.Text()
		return nil
	})

	wg.Wait()

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

	ctx.Send(
		fmt.Sprintf("Revise item created: \nTitle: %s\nDescription: %s", reviseItem.Name, reviseItem.Description),
		&tb.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]tb.InlineButton{
				{
					ReviseMenuButtonInline,
					ReviseCreateButtonInline,
				},
			},
		},
	)

	return nil
}
