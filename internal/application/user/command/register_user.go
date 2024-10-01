package command

import (
	"context"

	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
)

type RegisterUser struct {
	ChatID   domainUser.TelegramID `json:"chat_id"`
	Language string                `json:"language"`
	Settings *domainUser.Settings  `json:"settings"`
}

type RegisterUserHandler struct {
	userRepo domainUser.Repository
}

func NewRegisterUserHandler(userRepo domainUser.Repository) RegisterUserHandler {
	return RegisterUserHandler{userRepo: userRepo}
}

func (r RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	var opts []domainUser.OptionFunc
	if cmd.Settings != nil {
		opts = append(opts, domainUser.WithSettings(*cmd.Settings))
	}

	user, err := domainUser.NewUser(domainUser.NewUserID(), cmd.ChatID, opts...)
	if err != nil {
		return err
	}

	return r.userRepo.SaveUser(ctx, user)
}
