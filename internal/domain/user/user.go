package user

import (
	"time"

	"github.com/gofrs/uuid"
)

type TelegramID int64

type User struct {
	ID        uuid.UUID
	ChatID    TelegramID
	CreatedAt time.Time
	UpdatedAt time.Time
}
