package app

import (
	"log/slog"

	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/internal/tgbot"
)

type App struct {
	bot *tgbot.Bot
}

func NewApp(cfg config.Config, logger *slog.Logger) *App {

	bot, err := tgbot.NewBot(cfg.Telegram, logger)
	if err != nil {
		panic(err)
	}

	return &App{
		bot: bot,
	}
}

func (a *App) Start() {
	a.bot.Start()
}

func (a *App) Stop() {
	a.bot.Stop()
}
