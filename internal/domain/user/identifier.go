package user

import (
	"errors"
	"strconv"

	"github.com/gofrs/uuid"
)

var (
	ErrInvalidIdentifier = errors.New("invalid identifier")
)

type Identifier interface {
	// ID returns the identifier.
	ID() string
	// IsValid returns true if the identifier is valid.
	IsValid() bool
}

type UUIDIdentifier struct {
	id uuid.UUID
}

func (i UUIDIdentifier) ID() string {
	return i.id.String()
}

func (i UUIDIdentifier) UUID() uuid.UUID {
	return i.id
}

func (i UUIDIdentifier) IsValid() bool {
	return i.id != uuid.Nil
}

func NewUUIDIdentifier(stringID string) (UUIDIdentifier, error) {
	// Parse the string into a UUID.
	id, err := uuid.FromString(stringID)
	if err != nil {
		return UUIDIdentifier{}, err
	}

	return UUIDIdentifier{id: id}, nil
}

type TelegramIDWrapper struct {
	telegramID TelegramID
}

func (i TelegramIDWrapper) ID() string {
	return strconv.FormatInt(int64(i.telegramID), 10)
}

func (i TelegramIDWrapper) IsValid() bool {
	return i.telegramID != 0
}

func (i TelegramIDWrapper) TelegramID() TelegramID {
	return i.telegramID
}

func NewTelegramIDWrapper(telegramID TelegramID) TelegramIDWrapper {
	return TelegramIDWrapper{telegramID: telegramID}
}
