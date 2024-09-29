package user

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/text/language"
)

var (
	ErrInvalidSettings = errors.New("invalid settings")
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

// Settings represents user settings.
type Settings struct {
	// ID is used to identify the settings in the database.
	ID           uuid.UUID
	Language     language.Tag
	ReminderTime ReminderTime
}

// DefaultSettings returns default user settings.
func DefaultSettings() Settings {
	return Settings{
		Language:     language.English,
		ReminderTime: DefaultReminderTime(),
	}
}

func (s Settings) IsValid() bool {
	return s.Language != language.Und && s.ReminderTime.IsValid()
}

// ReminderTime represents a time of day when a reminder should be sent.
//
//		Note: support only 24-hour format. For example, 7:00 PM is represented as 19:00.
//		Maybe in the future, we will add support for 12-hour format. But not now.
//	 @@TODO: Add support for 12-hour format.
type ReminderTime struct {
	// uint8 is used to save space in the database
	Hour   uint8 // 0-23
	Minute uint8 // 0-59
}

func (r *ReminderTime) IsValid() bool {
	if r == nil {
		return false
	}
	return r.Hour <= 23 && r.Minute <= 59
}

// DefaultReminderTime returns a default reminder time.
func DefaultReminderTime() ReminderTime {
	return ReminderTime{
		Hour:   7,
		Minute: 0,
	}
}

// NewReminderTime creates a new reminder time.
//
//	hour maps to the 24-hour format, and minute is between 0-59.
func NewReminderTime(hour, minute uint8) ReminderTime {
	switch {
	case hour > 23:
		hour = 23
	case minute > 59:
		minute = 59
	}
	return ReminderTime{
		Hour:   hour,
		Minute: minute,
	}
}
