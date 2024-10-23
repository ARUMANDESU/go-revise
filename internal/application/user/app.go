package user

import (
	"github.com/ARUMANDESU/go-revise/internal/application/user/command"
	"github.com/ARUMANDESU/go-revise/internal/application/user/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterUser   command.RegisterUserHandler
	ChangeSettings command.ChangeSettingsHandler
}

type Queries struct {
	GetUser query.GetUserHandler
}
