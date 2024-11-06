package tgbot

import (
	"log/slog"
	"net/http"

	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/internal/application"
	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/internal/ports/tgbot/handler"
	"github.com/ARUMANDESU/go-revise/pkg/env"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
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
		// ParseMode: tb.ModeMarkdownV2, // WARNING: this will break the bot
		OnError: func(err error, c tb.Context) {
			if c != nil {
				slog.Error("error captured", logutil.Err(err), slog.Int("update_id", c.Update().ID))
			} else {
				slog.Error("error captured, context is nil", logutil.Err(err))
			}
		},
		Client: &httpClient,
		Poller: &webhookPoller,
	})
	if err != nil {
		return Port{}, err
	}

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
	p.bot.Start() // NOTE: this is blocking call
	return nil
}

func (p *Port) Stop() error {
	p.bot.Stop()
	return nil
}
