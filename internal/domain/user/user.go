package user

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

var (
	ErrInvalidUserID     = errors.New("invalid userID")
	ErrInvalidChatID     = errors.New("invalid chatID")
	ErrInvalidIdentifier = errors.New("invalid identifier")
)

// OptionFunc is a function that applies an option to a user.
// It is used to configure a user during creation.
type OptionFunc func(*User) error

type TelegramID int64

func (t TelegramID) IsValid() bool {
	return t != 0
}

func NewUserID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

// User represents a user, lol.
type User struct {
	id        uuid.UUID
	chatID    TelegramID
	createdAt time.Time
	updatedAt time.Time
	settings  Settings
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) ChatID() TelegramID {
	return u.chatID
}

func (u *User) Settings() Settings {
	return u.settings
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) UpdateSettings(settings Settings) error {
	op := errs.Op("domain.user.update_settings")
	if err := settings.Validate(); err != nil {
		return errs.WithOp(op, err, "invalid settings provided")
	}

	u.settings = settings
	u.updatedAt = time.Now()

	return nil
}

func NewUser(uid uuid.UUID, chatID TelegramID, options ...OptionFunc) (*User, error) {
	op := errs.Op("domain.user.new_user")
	switch {
	case uid == uuid.Nil:
		return nil, errs.
			NewIncorrectInputError(op, ErrInvalidUserID, "invalid userID").
			WithMessages([]errs.Message{{"message", "userID cannot be empty"}}).
			WithContext("userID", uid)
	}
	if !chatID.IsValid() {
		return nil, errs.
			NewIncorrectInputError(op, ErrInvalidChatID, "invalid chatID").
			WithMessages([]errs.Message{{"message", "chatID cannot be empty"}}).
			WithContext("chatID", chatID)
	}

	u := User{
		id:        uid,
		chatID:    chatID,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		settings:  DefaultSettings(),
	}

	for _, option := range options {
		if err := option(&u); err != nil {
			return nil, err
		}
	}

	return &u, nil
}

// MustNewUser creates a new user and panics if an error occurs.
//
//	Note: This function is intended for use in tests.
func MustNewUser(uid uuid.UUID, chatID TelegramID, options ...OptionFunc) *User {
	u, err := NewUser(uid, chatID, options...)
	if err != nil {
		panic(err)
	}
	return u
}

func WithSettings(settings Settings) OptionFunc {
	op := errs.Op("domain.user.with_settings")
	return func(u *User) error {
		if err := settings.Validate(); err != nil {
			return errs.WithOp(op, err, "invalid settings provided")
		}
		u.settings = settings
		return nil
	}
}

func WithCreatedAt(t time.Time) OptionFunc {
	return func(u *User) error {
		u.createdAt = t
		return nil
	}
}

func WithUpdatedAt(t time.Time) OptionFunc {
	return func(u *User) error {
		u.updatedAt = t
		return nil
	}
}
