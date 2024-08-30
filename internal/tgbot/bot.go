package tgbot

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/internal/domain"
	tb "gopkg.in/telebot.v3"
)

type ReviseService interface {
	Get(ctx context.Context, id string) (domain.ReviseItem, error)
	List(ctx context.Context, dto domain.ListReviseItemDTO) ([]domain.ReviseItem, domain.PaginationMetadata, error)
	Create(ctx context.Context, dto domain.CreateReviseItemDTO) (domain.ReviseItem, error)
	Update(ctx context.Context, dto domain.UpdateReviseItemDTO) (domain.ReviseItem, error)
	Delete(ctx context.Context, id string, userID string) (domain.ReviseItem, error)
}

type Bot struct {
	cfg        config.Telegram
	bot        *tb.Bot
	httpServer *http.Server
}

func NewBot(cfg config.Telegram, logger *slog.Logger) (*Bot, error) {

	httpServer := NewHTTPServer(cfg)

	webhook := &tb.Webhook{
		Listen:   cfg.URL,
		Endpoint: &tb.WebhookEndpoint{PublicURL: cfg.WebhookURL},
	}

	spamProtected := tb.NewMiddlewarePoller(webhook, func(upd *tb.Update) bool {
		if upd.Message == nil {
			return true
		}

		if strings.Contains(upd.Message.Text, "spam") {
			return false
		}

		return true
	})

	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: spamProtected,
	})
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:        b,
		cfg:        cfg,
		httpServer: httpServer,
	}, nil
}

func NewHTTPServer(cfg config.Telegram) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	httpServer := &http.Server{
		Addr:    cfg.WebhookURL,
		Handler: mux,
	}

	go httpServer.ListenAndServe()

	return httpServer
}

func (b *Bot) AddHandlers() {
	b.bot.Handle("/start", b.handleStartCommand)
	b.bot.Handle("/help", b.handleHelpCommand)

}

func (b *Bot) Start() error {
	b.AddHandlers()
	b.bot.Start()

	return nil
}

func (b *Bot) Stop() {
	b.bot.Stop()
	b.httpServer.Close()
}
