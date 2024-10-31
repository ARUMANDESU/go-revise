package user

import (
	"errors"
	"fmt"

	"golang.org/x/text/language"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

var ErrInvalidSettings = errors.New("invalid settings")

// Settings represents user settings.
type Settings struct {
	Language     language.Tag
	ReminderTime ReminderTime
}

func NewSettings(lang *language.Tag, reminderTime ReminderTime) (Settings, error) {
	switch {
	case lang == nil || *lang == language.Und:
		return Settings{}, errors.New("language is required")
	case !reminderTime.IsValid():
		return Settings{}, errors.New("reminder time is required")
	}

	return Settings{
		Language:     *lang,
		ReminderTime: reminderTime,
	}, nil
}

// DefaultSettings returns default user settings.
func DefaultSettings() Settings {
	return Settings{
		Language:     DefaultLanguage(),
		ReminderTime: DefaultReminderTime(),
	}
}

func (s *Settings) IsValid() bool {
	if s == nil {
		return false
	}

	return s.Language != language.Und && s.ReminderTime.IsValid()
}

func DefaultLanguage() language.Tag {
	return language.English
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
func NewReminderTime(hour, minute uint8) (ReminderTime, error) {
	const op = "domain.user.reminder_time.new_reminder_time"
	switch {
	case hour > 23:
		return ReminderTime{}, errs.NewIncorrectInputError(
			op,
			fmt.Errorf("hour must be less than 24(can be 23) and more than 0, got: %d", hour),
			"invalid hour",
		)
	case minute > 59:
		return ReminderTime{}, errs.NewIncorrectInputError(
			op,
			fmt.Errorf("minute must be less than 60(can be 59) and more than 0, got: %d", minute),
			"invalid minute",
		)
	}
	return ReminderTime{
		Hour:   hour,
		Minute: minute,
	}, nil
}
