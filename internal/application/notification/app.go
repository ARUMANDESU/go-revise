package notification

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/domain/reviseitem"
	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/retry"
)

// UserProvider defines the methods for user data access.
type UserProvider interface {
	// GetUsersForNotification selects users whose notify time is less than now and greater than now - 1 minute.
	GetUsersForNotification(ctx context.Context) ([]domainUser.User, error)
}

type ReviseItemProvider interface {
	FetchReviseItemsDueForUser(ctx context.Context, userID uuid.UUID) ([]reviseitem.Aggregate, error)
}

type Notifier interface {
	Notify(ctx context.Context, user domainUser.User, reviseItems []reviseitem.Aggregate) error
}

type Application struct {
	UserProvider       UserProvider
	ReviseItemProvider ReviseItemProvider
	Notifier           Notifier
}

func NewApplication(userProvider UserProvider, reviseItemProvider ReviseItemProvider, notifier Notifier) Application {
	return Application{
		UserProvider:       userProvider,
		ReviseItemProvider: reviseItemProvider,
		Notifier:           notifier,
	}
}

func (a Application) NotifyUsers(ctx context.Context) error {
	users, err := a.UserProvider.GetUsersForNotification(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users for notification: %w", err)
	}

	for _, user := range users {
		var reviseItems []reviseitem.Aggregate
		err = retry.Do(func() error {
			var err error
			reviseItems, err = a.ReviseItemProvider.FetchReviseItemsDueForUser(ctx, user.ID())
			if err != nil {
				return fmt.Errorf("failed to fetch revise items for user %s: %w", user.ID(), err)
			}
			return nil
		}, retry.WithMaxRetries(3))
		if err != nil {
			return err
		}

		err = retry.Do(func() error {
			if err := a.Notifier.Notify(ctx, user, reviseItems); err != nil {
				return fmt.Errorf("failed to notify user %s: %w", user.ID(), err)
			}
			return nil
		}, retry.WithMaxRetries(3))
	}

	return nil
}
