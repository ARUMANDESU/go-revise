package application_test

import (
	"context"
	"testing"

	"github.com/clarify/subtest"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"

	"github.com/ARUMANDESU/go-revise/internal/application"
	"github.com/ARUMANDESU/go-revise/internal/application/mocks"
	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
	"github.com/ARUMANDESU/go-revise/pkg/pointers"
)

type UserServiceSuite struct {
	userService        application.UserService
	mockUserProvider   *mocks.UserProvider
	mockUserRepository *mocks.UserRepository
}

func newUserServiceSuite(t *testing.T) UserServiceSuite {
	t.Helper()

	mockUserProvider := mocks.NewUserProvider(t)
	mockUserRepository := mocks.NewUserRepository(t)

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

func TestUserService_SaveUser(t *testing.T) {
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
				suite.mockUserRepository.On("Save", mock.Anything, mock.AnythingOfType("user.User")).Return(nil)
			},
			expectedError:  nil,
			expectedCalled: []string{"Save"},
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

			err := suite.userService.SaveUser(context.Background(), tt.params)

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
