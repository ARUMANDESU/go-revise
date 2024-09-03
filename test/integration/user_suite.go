package integration

import (
	"log/slog"
	"testing"

	usersvc "github.com/ARUMANDESU/go-revise/internal/service/user"
	"github.com/thejerf/slogassert"
)

type UserSuite struct {
	LogHandler *slogassert.Handler
	Service    usersvc.Service
}

func NewUserSuite(t *testing.T) *UserSuite {
	t.Helper()

	handler := slogassert.New(t, slog.LevelWarn, nil)
	log := slog.New(handler)

	storage, cleanup := setupSqlite(t)

	t.Cleanup(cleanup)

	return &UserSuite{
		LogHandler: handler,
		Service:    usersvc.NewService(log, storage, storage),
	}
}
