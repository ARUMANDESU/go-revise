package tgbot

import (
	"context"
	"fmt"
	"github.com/ARUMANDESU/go-revise/pkg/logger"
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

func (b *Bot) SendMessage(chatID int64, reviseItem domain.ReviseItem) error {
	const op = "tgbot.Bot.SendMessage"
	log := b.log.With("op", op)

	var message strings.Builder
	message.WriteString("You have a review item to revise\n")
	message.WriteString("Name: " + reviseItem.Name + "\n")
	message.WriteString("Description: " + reviseItem.Description + "\n")
	message.WriteString("Tags: " + strings.Join(reviseItem.Tags, ", ") + "\n")
	message.WriteString(fmt.Sprintf("Iteration: %d\n", reviseItem.Iteration))

	prevButton := tb.InlineButton{Text: domain.IntervalStringMap[domain.IntervalMap[reviseItem.Iteration-1]], Unique: "prev"}
	nextButton := tb.InlineButton{Text: domain.IntervalStringMap[domain.IntervalMap[reviseItem.Iteration+1]], Unique: "next"}

	_, err := b.bot.Send(
		&tb.Chat{ID: chatID},
		message.String(),
		&tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{
			{
				ResetButton,
			},
			{
				prevButton,
				nextButton,
			},
		}},
	)
	if err != nil {
		log.Error("failed to send message", logger.Err(err))
		return err
	}

	return nil
}
