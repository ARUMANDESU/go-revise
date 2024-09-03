package usersvc

import (
	"context"
	"errors"
	"testing"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/ARUMANDESU/go-revise/internal/service/user/mocks"
	"github.com/ARUMANDESU/go-revise/internal/storage"
	"github.com/ARUMANDESU/go-revise/pkg/logger"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type suite struct {
	service          Service
	mockUserProvider *mocks.UserProvider
	mockUserCreator  *mocks.UserCreator
}

func NewSuite(t *testing.T) *suite {
	t.Helper()

	mockUserProvider := mocks.NewUserProvider(t)
	mockUserCreator := mocks.NewUserCreator(t)
	service := NewService(logger.Plug(), mockUserProvider, mockUserCreator)
	return &suite{
		service:          service,
		mockUserProvider: mockUserProvider,
		mockUserCreator:  mockUserCreator,
	}
}

func TestService_GetUser(t *testing.T) {
	uid := uuid.Must(uuid.NewV7())

	tests := []struct {
		name          string
		userID        string
		expectedUser  domain.User
		expectedError error
		mockSetup     func(s *suite)
	}{
		{
			name:         "Happy Path",
			userID:       uid.String(),
			expectedUser: domain.User{ID: uid},
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUser", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{ID: uid}, nil)
			},
		},
		{
			name:          "error invalid argument",
			userID:        "",
			expectedError: service.ErrInvalidArgument,
			mockSetup:     func(s *suite) {},
		},
		{
			name:          "error invalid argument: invalid uuid",
			userID:        "invalid-uuid",
			expectedError: service.ErrInvalidArgument,
			mockSetup:     func(s *suite) {},
		},
		{
			name:          "error user not found",
			userID:        uid.String(),
			expectedError: service.ErrNotFound,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUser", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, storage.ErrNotFound)
			},
		},
		{
			name:          "unexpected error",
			userID:        uid.String(),
			expectedError: service.ErrInternal,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUser", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, errors.New("unexpected error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)
			ctx := context.Background()

			tt.mockSetup(s)
			user, err := s.service.Get(ctx, tt.userID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
			s.mockUserProvider.AssertExpectations(t)
		})
	}
}

func TestService_GetUserByChatID(t *testing.T) {

	var uid int64 = 12345

	tests := []struct {
		name          string
		chatID        int64
		expectedUser  domain.User
		expectedError error
		mockSetup     func(s *suite)
	}{
		{
			name:         "Happy Path",
			chatID:       uid,
			expectedUser: domain.User{TelegramID: uid},
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{TelegramID: uid}, nil)
			},
		},
		{
			name:          "user not found",
			chatID:        uid,
			expectedError: service.ErrNotFound,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{}, storage.ErrNotFound)
			},
		},
		{
			name:          "unexpected error",
			chatID:        uid,
			expectedError: service.ErrInternal,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{}, errors.New("unexpected error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)
			ctx := context.Background()

			tt.mockSetup(s)

			user, err := s.service.GetByChatID(ctx, tt.chatID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
			s.mockUserProvider.AssertExpectations(t)
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	var uid int64 = 12345

	tests := []struct {
		name          string
		chatID        int64
		expectedUser  domain.User
		expectedError error
		mockSetup     func(s *suite)
	}{
		{
			name:         "Happy Path",
			chatID:       uid,
			expectedUser: domain.User{TelegramID: uid},
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{}, storage.ErrNotFound)
				s.mockUserCreator.On("CreateUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(nil)
			},
		},
		{
			name:          "user already exists",
			chatID:        uid,
			expectedError: service.ErrAlreadyExists,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{ID: uuid.Must(uuid.NewV7()), TelegramID: uid}, nil)
			},
		},
		{
			name:          "unique constraint violation: on get user",
			chatID:        uid,
			expectedError: service.ErrAlreadyExists,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{}, storage.ErrAlreadyExists)
			},
		},
		{
			name:          "unique constraint violation: on create",
			chatID:        uid,
			expectedError: service.ErrAlreadyExists,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{}, storage.ErrNotFound)
				s.mockUserCreator.On("CreateUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(storage.ErrAlreadyExists)
			},
		},
		{
			name:          "unexpected error: failed to get user",
			chatID:        uid,
			expectedError: service.ErrInternal,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{}, errors.New("unexpected error"))
			},
		},
		{
			name:          "unexpected error: failed to create user",
			chatID:        uid,
			expectedError: service.ErrInternal,
			mockSetup: func(s *suite) {
				s.mockUserProvider.On("GetUserByChatID", mock.Anything, uid).Return(domain.User{}, storage.ErrNotFound)
				s.mockUserCreator.On("CreateUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(errors.New("unexpected error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)
			ctx := context.Background()

			tt.mockSetup(s)
			user, err := s.service.Create(ctx, tt.chatID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, domain.User{}, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.TelegramID, user.TelegramID)
			}
			s.mockUserCreator.AssertExpectations(t)
		})
	}
}
