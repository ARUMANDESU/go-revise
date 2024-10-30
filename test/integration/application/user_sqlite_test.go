package application

import (
	"testing"

	userapp "github.com/ARUMANDESU/go-revise/internal/application/user"
	usercommand "github.com/ARUMANDESU/go-revise/internal/application/user/command"
	userquery "github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/user/repository"
	"github.com/ARUMANDESU/go-revise/test/integration/tester"
)

func TestUserApp_RegisterUser(t *testing.T) {
}

func NewUserApplication(t *testing.T) userapp.Application {
	t.Helper()

	db := tester.NewSQLiteDB(t)
	userRepo := repository.NewSQLiteRepo(db)

	return userapp.Application{
		Commands: userapp.Commands{
			RegisterUser:   usercommand.NewRegisterUserHandler(&userRepo),
			ChangeSettings: usercommand.NewChangeSettingsHandler(&userRepo, &userRepo),
		},
		Queries: userapp.Queries{
			GetUser: userquery.NewGetUserHandler(&userRepo),
		},
	}
}
