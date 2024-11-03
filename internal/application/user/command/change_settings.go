package command

import (
	"context"

	"github.com/gofrs/uuid"

	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type UserProvider interface {
	GetUserByTelegramID(ctx context.Context, id domainUser.TelegramID) (*domainUser.User, error)
}

// ChangeSettings represents a command to change user settings.
// It can be used to change user settings by ID and chatID.
type ChangeSettings struct {
	ID       uuid.UUID             `json:"user_id"`
	ChatID   domainUser.TelegramID `json:"chat_id"`
	Settings domainUser.Settings   `json:"settings"`
}

type ChangeSettingsHandler struct {
	userRepo     domainUser.Repository
	userProvider UserProvider
}

func NewChangeSettingsHandler(userRepo domainUser.Repository, userProvider UserProvider) ChangeSettingsHandler {
	return ChangeSettingsHandler{
		userRepo:     userRepo,
		userProvider: userProvider,
	}
}

func (r ChangeSettingsHandler) Handle(ctx context.Context, cmd ChangeSettings) error {
	op := errs.Op("application.user.command.change_settings")
	if cmd.ID == uuid.Nil && !cmd.ChatID.IsValid() {
		return errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "id or chat ID must be provided").
			WithMessages([]errs.Message{{Key: "message", Value: "id or chat ID must be provided"}}).
			WithContext("chat_id", cmd.ChatID).
			WithContext("user_id", cmd.ID)
	}

	if cmd.ID == uuid.Nil && cmd.ChatID.IsValid() {
		user, err := r.userProvider.GetUserByTelegramID(ctx, cmd.ChatID)
		if err != nil {
			return errs.WithOp(op, err, "failed to get user by chat ID")
		}
		cmd.ID = user.ID()
	}

	err := r.userRepo.UpdateUser(ctx, cmd.ID, func(user *domainUser.User) (*domainUser.User, error) {
		if err := user.UpdateSettings(cmd.Settings); err != nil {
			return nil, errs.WithOp(op, err, "failed to update user settings")
		}
		return user, nil
	})
	if err != nil {
		return errs.WithOp(op, err, "failed to update user")
	}

	return nil
}
