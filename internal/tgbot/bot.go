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

type UserService interface {
	Get(ctx context.Context, id string) (domain.User, error)
	GetByChatID(ctx context.Context, chatID int64) (domain.User, error)
	Create(ctx context.Context, chatID int64) (domain.User, error)
}

type Bot struct {
	cfg           config.Telegram
	log           *slog.Logger
	bot           *tb.Bot
	httpServer    *http.Server
	ReviseService ReviseService
	UserService   UserService
}

func NewBot(cfg config.Telegram, logger *slog.Logger, reviseService ReviseService, userService UserService) (*Bot, error) {

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
		bot:           b,
		cfg:           cfg,
		log:           logger,
		httpServer:    httpServer,
		ReviseService: reviseService,
		UserService:   userService,
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

	b.bot.Handle(&ReviseMenuButtonInline, b.handleReviseMenuCommand)
	b.bot.Handle("/revise_menu", b.handleReviseMenuCommand)

	b.bot.Handle(&ReviseListButtonInline, b.handleReviseListCommand)
	b.bot.Handle("/revise_list", b.handleReviseListCommand)

	b.bot.Handle(&ReviseCreateButtonInline, b.handleReviseCreateCommand)
	b.bot.Handle("/revise_create", b.handleReviseCreateCommand)

	b.bot.SetCommands([]tb.Command{
		{Text: "start", Description: "Start the bot"},
		{Text: "help", Description: "Show help message"},
		{Text: "revise_menu", Description: "Revise commands"},
		{Text: "revise_list", Description: "List all revise items"},
		{Text: "revise_create", Description: "Create a new revise item"},
	})
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
