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
	// GetID returns the identifier as an interface{} to allow for both UUID and TelegramID.
	GetID() interface{}
	// GetIDAsString returns the identifier as a string.
	GetIDAsString() string
	// IsValid checks if the identifier is valid.
	IsValid() bool
}

type UUID uuid.UUID

func (u UUID) GetID() interface{} {
	return u
}

func (u UUID) GetIDAsString() string {
	return uuid.UUID(u).String()
}

func (u UUID) IsValid() bool {
	return uuid.UUID(u) != uuid.Nil
}

type TelegramID int64

func (t TelegramID) GetID() interface{} {
	return t
}

func (t TelegramID) GetIDAsString() string {
	return strconv.FormatInt(int64(t), 10)
}

func (t TelegramID) IsValid() bool {
	return t != 0
}
