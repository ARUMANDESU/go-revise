package tgbot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/internal/application"
	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/internal/ports/tgbot/handler"
	"github.com/ARUMANDESU/go-revise/internal/ports/tgbot/tgboterr"
	"github.com/ARUMANDESU/go-revise/pkg/env"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type Port struct {
	bot           *tb.Bot
	webhookPoller *tb.Webhook
	httpClient    *http.Client
	handler       *handler.Handler
}

func NewPort(cfg config.Telegram, app application.Application) (Port, error) {
	httpClient := http.Client{}

	webhookPoller := tb.Webhook{
		Listen: cfg.URL,
		Endpoint: &tb.WebhookEndpoint{
			PublicURL: cfg.WebhookURL,
		},
	}
	bot, err := tb.NewBot(tb.Settings{
		Token:   cfg.Token,
		Verbose: cfg.EnvMode != env.Prod,
		Offline: cfg.EnvMode == env.Test,
		OnError: tgboterr.OnError,
		Client:  &httpClient,
		Poller:  &webhookPoller,
	})
	if err != nil {
		return Port{}, err
	}

	slog.Info(
		"tgbot created",
		slog.String("url", cfg.URL),
		slog.String("webhook_url", cfg.WebhookURL),
	)

	return Port{
		bot:           bot,
		webhookPoller: &webhookPoller,
		httpClient:    &httpClient,
		handler:       handler.NewHandler(app),
	}, nil
}

// Start starts the bot.
// NOTE: this is blocking call
func (p *Port) Start() error {
	p.setUpRouter()
	slog.Info("tgbot is starting")
	p.bot.Start() // NOTE: this is blocking call
	return nil
}

func (p *Port) Stop() error {
	p.bot.Stop()
	slog.Info("tgbot is stopped")
	return nil
}

func (p *Port) Notify(
	ctx context.Context,
	user domainUser.User,
	reviseItems []reviseitem.ReviseItem,
) error {
	op := errs.Op("tgbot.port.notify")
	chat, err := p.bot.ChatByID(int64(user.ChatID()))
	if err != nil {
		return errs.WithOp(op, err, "failed to get chat").WithContext("chat_id", user.ChatID())
	}

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("Hello, %s!\n", chat.FirstName))
	msg.WriteString("You have the following revise items due:\n")

	_, err = p.bot.Send(tb.ChatID(user.ChatID()), msg.String())
	if err != nil {
		return errs.WithOp(op, err, "failed to notify user")
	}

	for _, reviseItem := range reviseItems {
		reviseItemMsg := strings.Builder{}
		reviseItemMsg.WriteString(fmt.Sprintf("Name: %s\n", reviseItem.Name()))
		if reviseItem.Description() != "" {
			reviseItemMsg.WriteString(fmt.Sprintf("Description: %s\n", reviseItem.Description()))
		}
		tags := reviseItem.Tags()
		if !tags.IsEmpty() {
			reviseItemMsg.WriteString(fmt.Sprintf("Tags: %s\n", tags.String()))
		}
		_, err = p.bot.Send(tb.ChatID(user.ChatID()), reviseItemMsg.String())
		if err != nil {
			return errs.WithOp(op, err, "failed to notify user")
		}
	}

	return nil
}
