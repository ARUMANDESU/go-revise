package user

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/pkg/i18n"
)

type TelegramID int64

type User struct {
	id        uuid.UUID
	chatID    TelegramID
	createdAt time.Time
	updatedAt time.Time
	settings  Settings
}

func WithSettings(settings Settings) func(*User) {
	return func(u *User) {
		u.settings = settings
	}
}

func NewUser(chatID TelegramID, options ...func(*User)) User {
	u := User{
		id:        uuid.Must(uuid.NewV4()),
		chatID:    chatID,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		settings:  DefaultSettings,
	}

	for _, option := range options {
		option(&u)
	}

	return u
}

func (u *User) UpdateSettings(settings Settings) {
	u.settings = settings
}

// Settings represents user settings.
type Settings struct {
	Language     i18n.Language
	ReminderTime ReminderTime
}

var DefaultSettings = Settings{
	Language:     i18n.DefaultLanguage,
	ReminderTime: DefaultReminderTime,
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

// DefaultReminderTime is a default time of day when a reminder should be sent.
// It is 7:00 AM.
var DefaultReminderTime = ReminderTime{
	Hour:   7,
	Minute: 0,
}
