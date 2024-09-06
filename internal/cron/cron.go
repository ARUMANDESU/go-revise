package cron

import (
	"context"
	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"time"
)

const ScanInterval = time.Minute

type TGNotifier interface {
	SendMessage(chatID int64, reviseItem domain.ReviseItem) error
}

type Scanner interface {
	Scan(ctx context.Context) ([]domain.ScheduledItem, error)
}

type Cron struct {
	log        *slog.Logger
	scheduler  gocron.Scheduler
	tgNotifier TGNotifier
	scanner    Scanner
}

func New(log *slog.Logger, notifier TGNotifier, scanner Scanner) (*Cron, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Cron{
		log:        log,
		scheduler:  s,
		tgNotifier: notifier,
		scanner:    scanner,
	}, nil
}

func (c *Cron) scanAndNotify() {
	c.log.Info("Scanning for items to notify")
	items, err := c.scanner.Scan(context.Background())
	if err != nil {
		// TODO: handle errors
		c.log.Error("failed to scan items", "error", err)
		return
	}

	for _, item := range items {
		job, err := c.scheduler.NewJob(
			gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(item.NotifyAt())),
			gocron.NewTask(func() {
				err := c.tgNotifier.SendMessage(item.TelegramID, item.ReviseItem)
				if err != nil {
					c.log.Error("failed to send message", "error", err)
				}
			}),
		)
		if err != nil {
			// TODO: handle error properly: keep track of failed jobs and retry them
			c.log.Error("failed to create job", "error", err)
			continue
		}

		err = job.RunNow()
		if err != nil {
			c.log.Error("failed to run job", "error", err)
		}
	}
}

func (c *Cron) Start() {
	const op = "cron.Cron.Start"
	log := c.log.With("op", op)

	log.Info("Starting cron job")
	job, err := c.scheduler.NewJob(
		gocron.DurationJob(ScanInterval),
		gocron.NewTask(c.scanAndNotify),
	)
	if err != nil {
		log.Error("failed to create cron job", "error", err)
	}

	c.scheduler.Start() // start is non-blocking, so we don't need to start it in a goroutine
	log.Info("cron jobs started")

	err = job.RunNow()
	if err != nil {
		log.Error("failed to run cron job", "error", err)
	}
}
