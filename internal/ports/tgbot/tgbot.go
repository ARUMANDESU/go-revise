package tgbot

import (
	"net/http"

	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/internal/application"
	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/internal/ports/tgbot/handler"
	"github.com/ARUMANDESU/go-revise/internal/ports/tgbot/tgboterr"
	"github.com/ARUMANDESU/go-revise/pkg/env"
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
