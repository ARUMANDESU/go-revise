package query

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/user"
)

type GetUserReadModel interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByChatID(ctx context.Context, chatID user.TelegramID) (User, error)
}

// GetUser represents a command to get a user.
// It can be used to get a user by ID or by chatID.
type GetUser struct {
	ID     uuid.UUID       `json:"user_id"`
	ChatID user.TelegramID `json:"chat_id"`
}

type GetUserHandler struct {
	readModel GetUserReadModel
}

func NewGetUserHandler(readModel GetUserReadModel) GetUserHandler {
	return GetUserHandler{readModel: readModel}
}

func (h GetUserHandler) Handle(ctx context.Context, cmd GetUser) (User, error) {
	if cmd.ID != uuid.Nil && !cmd.ChatID.IsValid() {
		return User{}, user.ErrInvalidChatID
	}

	if cmd.ID != uuid.Nil {
		return h.readModel.GetUserByID(ctx, cmd.ID)
	} else if cmd.ChatID.IsValid() {
		return h.readModel.GetUserByChatID(ctx, cmd.ChatID)
	}

	return User{}, user.ErrInvalidIdentifier
}
