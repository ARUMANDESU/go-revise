package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	adapterdb "github.com/ARUMANDESU/go-revise/internal/adapters/db"
	"github.com/ARUMANDESU/go-revise/internal/application"
	"github.com/ARUMANDESU/go-revise/internal/application/notification"
	reviseitemapp "github.com/ARUMANDESU/go-revise/internal/application/reviseitem"
	reviseitemcmd "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/command"
	reviseitemquery "github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
	userapp "github.com/ARUMANDESU/go-revise/internal/application/user"
	usercmd "github.com/ARUMANDESU/go-revise/internal/application/user/command"
	userquery "github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/domain/user/repository"
	httport "github.com/ARUMANDESU/go-revise/internal/ports/http"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}

	cfg := config.MustLoad()

	log, teardown := logutil.Setup(cfg.Env)
	defer teardown()

	log.Info("starting the app", slog.Attr{Key: "env", Value: slog.StringValue(cfg.Env.String())})

	db, err := adapterdb.NewSqlite(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to create new sqlite db", logutil.Err(err))
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Error("failed to ping db", logutil.Err(err))
		panic(err)
	}

	err = adapterdb.MigrateSchema(cfg.DatabaseURL, adapterdb.MigrationsFS, "", nil)
	if err != nil {
		log.Error("failed to migrate init", logutil.Err(err))
		panic(err)
	}

	userRepo := repository.NewSQLiteRepo(db)
	reviseitemRepo := reviseitem.NewSQLiteRepo(db)

	app := application.Application{
		User: userapp.Application{
			Commands: userapp.Commands{
				RegisterUser:   usercmd.NewRegisterUserHandler(&userRepo),
				ChangeSettings: usercmd.NewChangeSettingsHandler(&userRepo, &userRepo),
			},
			Queries: userapp.Queries{
				GetUser: userquery.NewGetUserHandler(&userRepo),
			},
		},
		ReviseItem: reviseitemapp.Application{
			Query: reviseitemapp.Query{
				GetReviseItem:       reviseitemquery.NewGetReviseItemHandler(&reviseitemRepo),
				ListUserReviseItems: reviseitemquery.NewListUserReviseItemsHandler(&reviseitemRepo),
			},
			Command: reviseitemapp.Command{
				NewReviseItem:     reviseitemcmd.NewNewReviseItemHandler(&reviseitemRepo),
				DeleteReviseItem:  reviseitemcmd.NewDeleteReviseItemHandler(&reviseitemRepo),
				ChangeDescription: reviseitemcmd.NewChangeDescriptionHandler(&reviseitemRepo),
				ChangeName:        reviseitemcmd.NewChangeNameHandler(&reviseitemRepo),
				AddTags:           reviseitemcmd.NewAddTagsHandler(&reviseitemRepo),
				RemoveTags:        reviseitemcmd.NewRemoveTagsHandler(&reviseitemRepo),
				Review:            reviseitemcmd.NewReviewHandler(&reviseitemRepo),
			},
		},
		Notification: notification.Application{}, // TODO: implement notification
	}

	httpPort := httport.NewHTTPPort(app)

	go func() {
		err := httpPort.Start(cfg.HTTP.Port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start http port", logutil.Err(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	defer log.Info("application stopped", slog.String("signal", sign.String()))
	log.Info("stopping application", slog.String("signal", sign.String()))

	err = httpPort.Stop()
}
