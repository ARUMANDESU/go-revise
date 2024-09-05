package app

import (
	"github.com/ARUMANDESU/go-revise/internal/cron"
	"log/slog"
	"sync"

	"github.com/ARUMANDESU/go-revise/internal/config"
	revisesvc "github.com/ARUMANDESU/go-revise/internal/service/revise"
	usersvc "github.com/ARUMANDESU/go-revise/internal/service/user"
	"github.com/ARUMANDESU/go-revise/internal/storage/sqlite"
	"github.com/ARUMANDESU/go-revise/internal/tgbot"
)

type App struct {
	wg        *sync.WaitGroup
	bot       *tgbot.Bot
	scheduler *cron.Cron
}

func NewApp(cfg config.Config, logger *slog.Logger) *App {

	var wg sync.WaitGroup

	sqliteDB, err := sqlite.NewStorage(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	reviseService := revisesvc.NewRevise(logger, &wg, revisesvc.ReviseStorages{ReviseProvider: sqliteDB, ReviseManager: sqliteDB, UserProvider: sqliteDB})
	userService := usersvc.NewService(logger, sqliteDB, sqliteDB)

	bot, err := tgbot.NewBot(cfg.Telegram, logger, &reviseService, &userService)
	if err != nil {
		panic(err)
	}

	scheduler, err := cron.New(logger, bot, reviseService)
	if err != nil {
		panic(err)
	}

	return &App{
		bot:       bot,
		scheduler: scheduler,
	}
}

func (a *App) Start() {
	go a.bot.Start()
	a.scheduler.Start()
}

func (a *App) Stop() {
	a.bot.Stop()
	// wait for all goroutines to finish
	a.wg.Wait()
}
