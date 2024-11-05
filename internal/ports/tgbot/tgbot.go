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

type TgBot struct {
	bot           *tb.Bot
	webhookPoller *tb.Webhook
	httpClient    *http.Client
	httpServer    *http.Server
	handler       *handler.Handler
}

func NewTgBot(cfg config.Telegram, app application.Application) (TgBot, error) {
	httpClient := http.Client{}
	httpServer := http.Server{
		Addr: cfg.WebhookURL,
	}

	go httpServer.ListenAndServe()

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
		// ParseMode: tb.ModeMarkdownV2,
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
		return TgBot{}, err
	}

	return TgBot{
		bot:           bot,
		webhookPoller: &webhookPoller,
		httpClient:    &httpClient,
		httpServer:    &httpServer,
		handler:       handler.NewHandler(app),
	}, nil
}

// Start starts the bot.
// NOTE: this is blocking call
func (t *TgBot) Start() error {
	t.setUpRouter()
	t.bot.Start() // NOTE: this is blocking call
	return nil
}

func (t *TgBot) Stop() error {
	t.bot.Stop()
	return t.httpServer.Close()
}
