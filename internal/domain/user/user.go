package user

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

// OptionFunc is a function that applies an option to a user.
// It is used to configure a user during creation.
type OptionFunc func(*User) error

// User represents a user, lol.
type User struct {
	id        uuid.UUID
	chatID    TelegramID
	createdAt time.Time
	updatedAt time.Time
	settings  Settings
}

func WithSettings(settings Settings) OptionFunc {
	return func(u *User) error {
		if !settings.IsValid() {
			return ErrInvalidSettings
		}
		u.settings = settings
		return nil
	}
}

func NewUser(chatID TelegramID, options ...OptionFunc) (User, error) {
	if !chatID.IsValid() {
		return User{}, errors.New("chatID is required")
	}

	u := User{
		id:        NewUserID(),
		chatID:    chatID,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		settings:  DefaultSettings(),
	}

	for _, option := range options {
		if err := option(&u); err != nil {
			return User{}, err
		}
	}

	return u, nil
}

// MustNewUser creates a new user and panics if an error occurs.
//
//	Note: This function is intended for use in tests.
func MustNewUser(chatID TelegramID, options ...OptionFunc) User {
	u, err := NewUser(chatID, options...)
	if err != nil {
		panic(err)
	}
	return u
}

func (u *User) UpdateSettings(settings Settings) {
	u.settings = settings
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func NewUserID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
