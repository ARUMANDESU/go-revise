package app

import (
	"log/slog"

	"github.com/ARUMANDESU/go-revise/internal/config"
	revisesvc "github.com/ARUMANDESU/go-revise/internal/service/revise"
	usersvc "github.com/ARUMANDESU/go-revise/internal/service/user"
	"github.com/ARUMANDESU/go-revise/internal/storage/sqlite"
	"github.com/ARUMANDESU/go-revise/internal/tgbot"
)

type App struct {
	bot *tgbot.Bot
}

func NewApp(cfg config.Config, logger *slog.Logger) *App {

	sqliteDB, err := sqlite.NewStorage(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	reviseService := revisesvc.NewRevise(logger, revisesvc.ReviseStorages{ReviseProvider: sqliteDB, ReviseManager: sqliteDB, UserProvider: sqliteDB})
	userService := usersvc.NewService(logger, sqliteDB, sqliteDB)

	bot, err := tgbot.NewBot(cfg.Telegram, logger, &reviseService, &userService)
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
