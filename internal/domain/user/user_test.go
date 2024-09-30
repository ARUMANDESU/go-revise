package user

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/ARUMANDESU/go-revise/pkg/pointers"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	userID := NewUserID()
	tests := []struct {
		name          string
		userID        uuid.UUID
		telegramID    TelegramID
		options       []OptionFunc
		want          User
		expectedError error
	}{
		{
			name:       "With valid data",
			userID:     userID,
			telegramID: 123456789,
			want: User{
				id:        userID,
				chatID:    123456789,
				createdAt: time.Now(),
				updatedAt: time.Now(),
				settings:  DefaultSettings(),
			},
			expectedError: nil,
		},
		{
			name:       "With valid custom settings",
			userID:     userID,
			telegramID: 123456789,
			options:    []OptionFunc{WithSettings(validSettings(t))},
			want: User{
				id:        userID,
				chatID:    123456789,
				createdAt: time.Now(),
				updatedAt: time.Now(),
				settings:  validSettings(t),
			},
		},
		{
			name:          "With invalid telegram ID",
			userID:        userID,
			telegramID:    0,
			want:          User{},
			expectedError: ErrInvalidChatID,
		},
		{
			name:          "With invalid user ID",
			userID:        uuid.Nil,
			telegramID:    123456789,
			want:          User{},
			expectedError: ErrInvalidUserID,
		},
		{
			name:          "With invalid custom settings",
			userID:        userID,
			telegramID:    123456789,
			options:       []OptionFunc{WithSettings(invalidSettings(t))},
			want:          User{},
			expectedError: ErrInvalidSettings,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.userID, tt.telegramID, tt.options...)

			t.Run("Expect error", func(t *testing.T) {
				require.ErrorIs(t, err, tt.expectedError)
			})
			if err == nil {
				t.Run("Expect user", func(t *testing.T) {
					AssertUser(t, got, tt.want)
				})
			}
		})
	}
}

func TestUser_UpdateSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		settings      Settings
		expectedError error
	}{
		{
			name:          "With valid settings",
			settings:      validSettings(t),
			expectedError: nil,
		},
		{
			name:          "With invalid settings",
			settings:      invalidSettings(t),
			expectedError: ErrInvalidSettings,
		},
		{
			name:          "With nil settings",
			settings:      Settings{},
			expectedError: ErrInvalidSettings,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := MustNewUser(NewUserID(), 123456789)
			err := user.UpdateSettings(tt.settings)

			t.Run("Expect error", func(t *testing.T) {
				require.ErrorIs(t, err, tt.expectedError)
			})
			if err == nil {
				t.Run("Expect settings", func(t *testing.T) {
					require.Equal(t, tt.settings, user.Settings())
					require.WithinDuration(t, user.updatedAt, time.Now(), time.Second)
				})
			}
		})
	}
}

func AssertUser(t *testing.T, got, want User) {
	t.Helper()

	require.Equal(t, want.id, got.id)
	require.Equal(t, want.chatID, got.chatID)
	require.WithinDuration(t, want.createdAt, got.createdAt, time.Second)
	require.WithinDuration(t, want.updatedAt, got.updatedAt, time.Second)
	require.Equal(t, want.settings, got.settings)
}

func validSettings(t *testing.T) Settings {
	t.Helper()

	s, err := NewSettings(pointers.New(language.Kazakh), DefaultReminderTime())
	if err != nil {
		t.Fatal(err)
	}

	return s
}

func invalidSettings(t *testing.T) Settings {
	t.Helper()

	return Settings{}
}
