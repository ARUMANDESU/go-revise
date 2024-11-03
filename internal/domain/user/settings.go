package user

import (
	"errors"

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
	op := errs.Op("domain.user.new_settings")
	if lang == nil || *lang == language.Und {
		return Settings{}, errs.
			NewIncorrectInputError(op, ErrInvalidSettings, "language is not provided").
			WithMessages([]errs.Message{{"message", "language is not provided"}}).
			WithContext("language", lang)
	}
	if err := reminderTime.Validate(); err != nil {
		return Settings{}, errs.WithOp(op, err, "reminder time is invalid")
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

func (s *Settings) Validate() error {
	op := errs.Op("domain.user.settings.validate")
	if s == nil {
		return errs.
			NewIncorrectInputError(op, ErrInvalidSettings, "settings are nil").
			WithMessages([]errs.Message{{"message", "settings is not provided"}})
	}
	if s.Language == language.Und {
		return errs.
			NewIncorrectInputError(op, ErrInvalidSettings, "language is not set").
			WithMessages([]errs.Message{{"message", "language is not set"}}).
			WithContext("language", s.Language)
	}
	if err := s.ReminderTime.Validate(); err != nil {
		return errs.
			NewIncorrectInputError(op, ErrInvalidSettings, "reminder time is invalid").
			WithMessages([]errs.Message{{"message", "reminder time is invalid"}}).
			WithContext("reminder_time", s.ReminderTime)
	}
	return nil
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

func (r *ReminderTime) Validate() error {
	const op = "domain.user.reminder_time.validate"
	if r == nil {
		return errs.
			NewIncorrectInputError(op, ErrInvalidSettings, "reminder time is nil").
			WithMessages([]errs.Message{{"message", "reminder time is not provided"}})
	}
	if r.Hour > 23 {
		return errs.
			NewIncorrectInputError(op, ErrInvalidSettings, "hour is invalid").
			WithMessages([]errs.Message{{"message", "hour is invalid"}}).
			WithContext("hour", r.Hour)
	}
	if r.Minute > 59 {
		return errs.
			NewIncorrectInputError(op, ErrInvalidSettings, "minute is invalid").
			WithMessages([]errs.Message{{"message", "minute is invalid"}}).
			WithContext("minute", r.Minute)
	}
	return nil
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
	const op = "domain.user.new_reminder_time"
	switch {
	case hour > 23:
		return ReminderTime{}, errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "invalid hour").
			WithMessages([]errs.Message{{"message", "hour must be less than 24 and more than 0"}}).
			WithContext("hour", hour)
	case minute > 59:
		return ReminderTime{}, errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "invalid minute").
			WithMessages([]errs.Message{{"message", "minute must be less than 60 and more than 0"}}).
			WithContext("minute", minute)
	}
	return ReminderTime{
		Hour:   hour,
		Minute: minute,
	}, nil
}
