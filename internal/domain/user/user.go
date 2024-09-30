package user

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrInvalidUserID = errors.New("invalid userID")
	ErrInvalidChatID = errors.New("invalid chatID")
)

// OptionFunc is a function that applies an option to a user.
// It is used to configure a user during creation.
type OptionFunc func(*User) error

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

func (u *User) UpdateSettings(settings Settings) error {
	if !settings.IsValid() {
		return ErrInvalidSettings
	}

	u.settings = settings
	u.updatedAt = time.Now()

	return nil
}

func (u *User) Settings() Settings {
	return u.settings
}

func NewUser(uid uuid.UUID, chatID TelegramID, options ...OptionFunc) (User, error) {
	switch {
	case uid == uuid.Nil:
		return User{}, ErrInvalidUserID
	}
	if !chatID.IsValid() {
		return User{}, ErrInvalidChatID
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
			return User{}, err
		}
	}

	return u, nil
}

// MustNewUser creates a new user and panics if an error occurs.
//
//	Note: This function is intended for use in tests.
func MustNewUser(uid uuid.UUID, chatID TelegramID, options ...OptionFunc) User {
	u, err := NewUser(uid, chatID, options...)
	if err != nil {
		panic(err)
	}
	return u
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
