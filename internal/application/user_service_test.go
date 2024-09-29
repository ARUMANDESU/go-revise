package application_test

import (
	"context"
	"testing"

	"github.com/clarify/subtest"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"

	"github.com/ARUMANDESU/go-revise/internal/application"
	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/internal/domain/user/mocks"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
	"github.com/ARUMANDESU/go-revise/pkg/pointers"
)

type UserServiceSuite struct {
	userService        application.UserService
	mockUserProvider   *mocks.Provider
	mockUserRepository *mocks.Repository
}

func newUserServiceSuite(t *testing.T) UserServiceSuite {
	t.Helper()

	mockUserProvider := mocks.NewProvider(t)
	mockUserRepository := mocks.NewRepository(t)

	userService := application.NewUserService(logutil.Plug(), mockUserRepository, mockUserProvider)

	return UserServiceSuite{
		userService:        userService,
		mockUserProvider:   mockUserProvider,
		mockUserRepository: mockUserRepository,
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		id             domainUser.Identifier
		mockSetup      func(suite UserServiceSuite)
		expectedError  error
		expectedCalled []string
	}{
		{
			name: "With uuid",
			id:   domainUser.NewUserUUID(),
			mockSetup: func(suite UserServiceSuite) {
				suite.mockUserProvider.On("GetUserByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(domainUser.User{}, nil)
			},
			expectedError:  nil,
			expectedCalled: []string{"GetUserByID"},
		},
		{
			name: "With telegram id",
			id:   domainUser.NewTelegramID(123456789),
			mockSetup: func(suite UserServiceSuite) {
				suite.mockUserProvider.On("GetUserByTelegramID", mock.Anything, mock.AnythingOfType("user.TelegramID")).Return(domainUser.User{}, nil)
			},
			expectedError:  nil,
			expectedCalled: []string{"GetUserByTelegramID"},
		},
		{
			name:           "With invalid identifier",
			id:             domainUser.TelegramID(0),
			mockSetup:      func(suite UserServiceSuite) {},
			expectedError:  domainUser.ErrInvalidIdentifier,
			expectedCalled: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newUserServiceSuite(t)
			tt.mockSetup(suite)

			user, err := suite.userService.GetUserByID(context.Background(), tt.id)

			for _, method := range tt.expectedCalled {
				suite.mockUserProvider.AssertCalled(t, method, mock.Anything, mock.Anything)
			}
			suite.mockUserProvider.AssertExpectations(t)
			suite.mockUserRepository.AssertExpectations(t)

			if tt.expectedError != nil {
				t.Run("Expect error", subtest.Value(err).ErrorIs(tt.expectedError))
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
				t.Run("Expect user", subtest.Value(user).NotReflectNil())
			}
		})
	}
}

func TestUserService_RegisterUser(t *testing.T) {
	tests := []struct {
		name           string
		params         application.NewUserServiceParams
		mockSetup      func(suite UserServiceSuite)
		expectedError  error
		expectedCalled []string
	}{
		{
			name: "With valid params",
			params: application.NewUserServiceParams{
				ChatID:       domainUser.NewTelegramID(123456789),
				Language:     language.Kazakh,
				ReminderTime: pointers.New(domainUser.NewReminderTime(10, 0)),
			},
			mockSetup: func(suite UserServiceSuite) {
				suite.mockUserRepository.On("SaveUser", mock.Anything, mock.AnythingOfType("user.User")).Return(nil)
			},
			expectedError:  nil,
			expectedCalled: []string{"SaveUser"},
		},
		{
			name: "With invalid params",
			params: application.NewUserServiceParams{
				ChatID:       domainUser.NewTelegramID(0),
				Language:     language.Kazakh,
				ReminderTime: nil,
			},
			mockSetup:      func(suite UserServiceSuite) {},
			expectedError:  application.ErrInvalidArguments,
			expectedCalled: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newUserServiceSuite(t)
			tt.mockSetup(suite)

			err := suite.userService.RegisterUser(context.Background(), tt.params)

			for _, method := range tt.expectedCalled {
				suite.mockUserRepository.AssertCalled(t, method, mock.Anything, mock.Anything)
			}
			suite.mockUserProvider.AssertExpectations(t)
			suite.mockUserRepository.AssertExpectations(t)

			if tt.expectedError != nil {
				t.Run("Expect error", subtest.Value(err).ErrorIs(tt.expectedError))
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
			}
		})
	}
}

func TestUserService_UpdateUserSettings(t *testing.T) {
	tests := []struct {
		name           string
		id             domainUser.Identifier
		settings       domainUser.Settings
		mockSetup      func(suite UserServiceSuite)
		expectedError  error
		expectedCalled []string
	}{
		{
			name:     "With valid params",
			id:       domainUser.NewUserUUID(),
			settings: domainUser.DefaultSettings(),
			mockSetup: func(suite UserServiceSuite) {
				suite.mockUserRepository.On(
					"UpdateUser",
					mock.Anything,
					mock.AnythingOfType("uuid.UUID"),
					mock.Anything, // updateFunc
				).Return(nil)
			},
			expectedError:  nil,
			expectedCalled: []string{"UpdateUser"},
		},
		{
			name:           "With invalid identifier",
			id:             domainUser.TelegramID(0),
			settings:       domainUser.DefaultSettings(),
			mockSetup:      func(suite UserServiceSuite) {},
			expectedError:  domainUser.ErrInvalidIdentifier,
			expectedCalled: []string{},
		},
		{
			name: "With invalid settings",
			id:   domainUser.NewUserUUID(),
			settings: domainUser.Settings{
				Language:     language.Und,
				ReminderTime: domainUser.ReminderTime{},
			},
			mockSetup:      func(suite UserServiceSuite) {},
			expectedError:  domainUser.ErrInvalidSettings,
			expectedCalled: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newUserServiceSuite(t)
			tt.mockSetup(suite)

			err := suite.userService.UpdateUserSettings(context.Background(), tt.id, tt.settings)

			for _, method := range tt.expectedCalled {
				suite.mockUserRepository.AssertCalled(t, method, mock.Anything, mock.Anything, mock.Anything)
			}
			suite.mockUserProvider.AssertExpectations(t)
			suite.mockUserRepository.AssertExpectations(t)

			if tt.expectedError != nil {
				t.Run("Expect error", subtest.Value(err).ErrorIs(tt.expectedError))
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
			}
		})
	}
}
