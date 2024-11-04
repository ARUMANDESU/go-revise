package command

import (
	"context"

	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type RegisterUser struct {
	ChatID   domainUser.TelegramID `json:"chat_id"`
	Settings *domainUser.Settings  `json:"settings"`
}

type RegisterUserHandler struct {
	userRepo domainUser.Repository
}

func NewRegisterUserHandler(userRepo domainUser.Repository) RegisterUserHandler {
	return RegisterUserHandler{userRepo: userRepo}
}

func (r RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	const op = "application.user.register_user"
	var opts []domainUser.OptionFunc
	if cmd.Settings != nil {
		opts = append(opts, domainUser.WithSettings(*cmd.Settings))
	}

	user, err := domainUser.NewUser(domainUser.NewUserID(), cmd.ChatID, opts...)
	if err != nil {
		return errs.WithOp(op, err, "failed to create user")
	}

	err = r.userRepo.CreateUser(ctx, *user)
	if err != nil {
		return errs.WithOp(op, err, "failed to create user")
	}
	return nil
}
