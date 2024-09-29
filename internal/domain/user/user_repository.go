package user

import (
	"context"

	"github.com/gofrs/uuid"
)

// Provider handles the retrieval of user data.
//
//go:generate mockery --name Provider --output mocks/
type Provider interface {
	// GetUserByID returns a user by ID(UUID).
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	// GetUserByTelegramID returns a user by Telegram ID(int64).
	GetUserByTelegramID(ctx context.Context, telegramID TelegramID) (User, error)
}

// Repository handles the persistence of user data.
//
//go:generate mockery --name Repository --output mocks/
type Repository interface {
	// SaveUser saves a user.
	SaveUser(ctx context.Context, u User) error
	// UpdateUser updates user data.
	// The updateFn function is called with the user data to be updated.
	UpdateUser(ctx context.Context, userID uuid.UUID, updateFn func(*User) (*User, error)) error
}
