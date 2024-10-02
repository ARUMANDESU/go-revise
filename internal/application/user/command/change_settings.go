package command

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
)

type UserProvider interface {
	GetUserByTelegramID(ctx context.Context, id domainUser.TelegramID) (domainUser.User, error)
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
	if cmd.ID != uuid.Nil && !cmd.ChatID.IsValid() {
		return fmt.Errorf("must provide either user_id or chat_id")
	}

	if cmd.ID == uuid.Nil && cmd.ChatID.IsValid() {
		user, err := r.userProvider.GetUserByTelegramID(ctx, cmd.ChatID)
		if err != nil {
			return err
		}
		cmd.ID = user.ID()
	}

	return r.userRepo.UpdateUser(ctx, cmd.ID, func(user *domainUser.User) (*domainUser.User, error) {
		err := user.UpdateSettings(cmd.Settings)
		if err != nil {
			return nil, err
		}
		return user, nil
	})
}