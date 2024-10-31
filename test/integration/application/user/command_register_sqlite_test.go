package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userapp "github.com/ARUMANDESU/go-revise/internal/application/user"
	usercommand "github.com/ARUMANDESU/go-revise/internal/application/user/command"
	userquery "github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/internal/domain/user/repository"
	"github.com/ARUMANDESU/go-revise/test/integration/tester"
)

func TestUserApp_RegisterUser(t *testing.T) {
	tests := []struct {
		name string
		cmd  usercommand.RegisterUser
	}{
		{
			name: "With valid command",
			cmd: usercommand.RegisterUser{
				ChatID: user.TelegramID(123),
				Settings: &user.Settings{
					Language:     user.DefaultLanguage(),
					ReminderTime: user.DefaultReminderTime(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userApp := NewUserApplication(t)

			err := userApp.Commands.RegisterUser.Handle(ctx, tt.cmd)
			require.NoError(t, err, "failed to register user")

			queryUser, err := userApp.Queries.GetUser.Handle(
				ctx,
				userquery.GetUser{ChatID: tt.cmd.ChatID},
			)
			require.NoError(t, err, "failed to get user")

			assert.Equal(t, tt.cmd.ChatID, user.TelegramID(queryUser.ChatID))
			assert.Equal(t, tt.cmd.Settings.ReminderTime.Hour, queryUser.Settings.ReminderTime.Hour)
			assert.Equal(
				t,
				tt.cmd.Settings.ReminderTime.Minute,
				queryUser.Settings.ReminderTime.Minute,
			)
		})
	}
}

func NewUserApplication(t *testing.T) userapp.Application {
	t.Helper()

	db := tester.NewSQLiteDB(t)
	userRepo := repository.NewSQLiteRepo(db)

	return userapp.Application{
		Commands: userapp.Commands{
			RegisterUser:   usercommand.NewRegisterUserHandler(&userRepo),
			ChangeSettings: usercommand.NewChangeSettingsHandler(&userRepo, &userRepo),
		},
		Queries: userapp.Queries{
			GetUser: userquery.NewGetUserHandler(&userRepo),
		},
	}
}
