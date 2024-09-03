package integration

import (
	"context"
	"testing"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_User_Get(t *testing.T) {
	ctx := context.Background()
	s := NewUserSuite(t)

	tests := []struct {
		name          string
		userID        string
		expectedError error
		expectedUser  domain.User
	}{
		{
			name:          "found",
			userID:        "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			expectedError: nil,
			expectedUser: domain.User{
				ID:         uuid.FromStringOrNil("1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e"),
				TelegramID: 123456789,
			},
		},
		{
			name:          "not found",
			userID:        "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b88",
			expectedError: service.ErrNotFound,
		},
		{
			name:          "invalid id",
			userID:        "invalid",
			expectedError: service.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := s.Service.Get(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				s.LogHandler.AssertEmpty()
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedUser.ID.String(), user.ID.String())
				assert.Equal(t, tt.expectedUser.TelegramID, user.TelegramID)
			}
		})
	}
}

func TestService_User_GetByChatID(t *testing.T) {
	ctx := context.Background()
	s := NewUserSuite(t)

	tests := []struct {
		name          string
		chatID        int64
		expectedError error
		expectedUser  domain.User
	}{
		{
			name:          "found",
			chatID:        123456789,
			expectedError: nil,
			expectedUser: domain.User{
				ID:         uuid.FromStringOrNil("1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e"),
				TelegramID: 123456789,
			},
		},
		{
			name:          "not found",
			chatID:        123456788,
			expectedError: service.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := s.Service.GetByChatID(ctx, tt.chatID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				s.LogHandler.AssertEmpty()
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedUser.ID.String(), user.ID.String())
				assert.Equal(t, tt.expectedUser.TelegramID, user.TelegramID)
			}
		})
	}
}

func TestService_User_Create(t *testing.T) {
	ctx := context.Background()
	s := NewUserSuite(t)

	tests := []struct {
		name          string
		chatID        int64
		expectedError error
	}{
		{
			name:          "created",
			chatID:        123456788,
			expectedError: nil,
		},
		{
			name:          "already exists",
			chatID:        123456789,
			expectedError: service.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := s.Service.Create(ctx, tt.chatID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				s.LogHandler.AssertEmpty()
				assert.NoError(t, err)

				assert.NotEqual(t, uuid.Nil, user.ID)
				assert.Equal(t, tt.chatID, user.TelegramID)

				// Check if the user is created

				user, err := s.Service.GetByChatID(ctx, tt.chatID)
				assert.NoError(t, err)

				assert.NotEqual(t, uuid.Nil, user.ID)
				assert.Equal(t, tt.chatID, user.TelegramID)
			}
		})
	}
}
