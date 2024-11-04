package query

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
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
	op := errs.Op("application.user.query.get_user")
	if cmd.ID != uuid.Nil && !cmd.ChatID.IsValid() {
		return User{}, errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "id or chat ID must be provided").
			WithMessages([]errs.Message{{Key: "message", Value: "id or chat ID must be provided"}}).
			WithContext("chat_id", cmd.ChatID).
			WithContext("user_id", cmd.ID)
	}

	if cmd.ID != uuid.Nil {
		queryUser, err := h.readModel.GetUserByID(ctx, cmd.ID)
		if err != nil {
			return User{}, errs.WithOp(op, err, "failed to get user by ID")
		}
		return queryUser, nil
	} else if cmd.ChatID.IsValid() {
		queryUser, err := h.readModel.GetUserByChatID(ctx, cmd.ChatID)
		if err != nil {
			return User{}, errs.WithOp(op, err, "failed to get user by chat ID")
		}
		return queryUser, nil
	}

	return User{}, user.ErrInvalidIdentifier
}
