package application_test

import (
	"context"
	"testing"

	"github.com/clarify/subtest"
	"github.com/stretchr/testify/mock"

	"github.com/ARUMANDESU/go-revise/internal/application"
	"github.com/ARUMANDESU/go-revise/internal/application/mocks"
	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
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

	t.Run("With uuid", func(t *testing.T) {
		suite := newUserServiceSuite(t)
		userID := domainUser.NewUserUUID()

		suite.mockUserProvider.On("GetUserByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(domainUser.User{}, nil)
		user, err := suite.userService.GetUserByID(context.Background(), userID)

		t.Run("Expect GetUserByID called", func(t *testing.T) {
			suite.mockUserProvider.AssertCalled(t, "GetUserByID", mock.Anything, mock.AnythingOfType("uuid.UUID"))
		})
		t.Run("Expect GetUserByTelegramID not called", func(t *testing.T) {
			suite.mockUserProvider.AssertNotCalled(t, "GetUserByTelegramID")
		})
		t.Run("Expect Save not called", func(t *testing.T) {
			suite.mockUserRepository.AssertNotCalled(t, "Save")
		})
		t.Run("Expect UpdateSettings not called", func(t *testing.T) {
			suite.mockUserRepository.AssertNotCalled(t, "UpdateSettings")
		})
		t.Run("Expect no error", subtest.Value(err).NoError())
		t.Run("Expect user", subtest.Value(user).NotReflectNil())
	})

	t.Run("With telegram id", func(t *testing.T) {
		suite := newUserServiceSuite(t)
		telegramID := domainUser.NewTelegramID(123456789)

		suite.mockUserProvider.On(
			"GetUserByTelegramID",
			mock.Anything,
			mock.AnythingOfType("user.TelegramID"),
		).Return(domainUser.User{}, nil)
		user, err := suite.userService.GetUserByID(context.Background(), telegramID)

		t.Run("Expect GetUserByID not called", func(t *testing.T) {
			suite.mockUserProvider.AssertNotCalled(t, "GetUserByID")
		})
		t.Run("Expect GetUserByTelegramID called", func(t *testing.T) {
			suite.mockUserProvider.AssertCalled(t, "GetUserByTelegramID", mock.Anything, mock.AnythingOfType("user.TelegramID"))
		})
		t.Run("Expect Save not called", func(t *testing.T) {
			suite.mockUserRepository.AssertNotCalled(t, "Save")
		})
		t.Run("Expect UpdateSettings not called", func(t *testing.T) {
			suite.mockUserRepository.AssertNotCalled(t, "UpdateSettings")
		})
		t.Run("Expect no error", subtest.Value(err).NoError())
		t.Run("Expect user", subtest.Value(user).NotReflectNil())
	})

	t.Run("With invalid identifier", func(t *testing.T) {
		suite := newUserServiceSuite(t)
		invalidID := domainUser.TelegramID(0)

		user, err := suite.userService.GetUserByID(context.Background(), invalidID)

		t.Run("Expect GetUserByID not called", func(t *testing.T) {
			suite.mockUserProvider.AssertNotCalled(t, "GetUserByID")
		})
		t.Run("Expect GetUserByTelegramID not called", func(t *testing.T) {
			suite.mockUserProvider.AssertNotCalled(t, "GetUserByTelegramID")
		})
		t.Run("Expect Save not called", func(t *testing.T) {
			suite.mockUserRepository.AssertNotCalled(t, "Save")
		})
		t.Run("Expect UpdateSettings not called", func(t *testing.T) {
			suite.mockUserRepository.AssertNotCalled(t, "UpdateSettings")
		})
		t.Run("Expect error", subtest.Value(err).ErrorIs(domainUser.ErrInvalidIdentifier))
		t.Run("Expect user", subtest.Value(user).DeepEqual(domainUser.User{}))
	})
}
