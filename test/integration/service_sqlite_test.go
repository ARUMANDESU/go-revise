package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Revise_Get(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name         string
		id           string
		expectedItem domain.ReviseItem
	}{
		{
			name: "existing item",
			id:   "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			expectedItem: domain.ReviseItem{
				ID:          uuid.FromStringOrNil("3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e"),
				UserID:      uuid.FromStringOrNil("1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e"),
				Name:        "First Revise Item",
				Description: "Description for first revise item",
				Tags:        []string{"tag1", "tag2"},
				Iteration:   0,
			},
		},
		{
			name: "existing item 2",
			id:   "4e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			expectedItem: domain.ReviseItem{
				ID:          uuid.FromStringOrNil("4e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e"),
				UserID:      uuid.FromStringOrNil("2e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e"),
				Name:        "Second Revise Item",
				Description: "Description for second revise item",
				Tags:        []string{"tag3", "tag4"},
				Iteration:   2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := s.Service.Get(ctx, tt.id)

			assert.NoError(t, err)

			assert.Equal(t, tt.expectedItem.ID.String(), item.ID.String())
			assert.Equal(t, tt.expectedItem.UserID.String(), item.UserID.String())
			assert.Equal(t, tt.expectedItem.Name, item.Name)
			assert.Equal(t, tt.expectedItem.Description, item.Description)
			assert.Equal(t, tt.expectedItem.Tags, item.Tags)
			assert.Equal(t, tt.expectedItem.Iteration, item.Iteration)
		})
	}
}

func TestService_Revise_Get_FailPath(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "non-existing item",
			id:            "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8f",
			expectedError: service.ErrNotFound,
		},
		{
			name:          "invalid id",
			id:            "invalid",
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "empty id",
			id:            "",
			expectedError: service.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Service.Get(ctx, tt.id)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}

}

func TestService_Revise_Create(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name         string
		dto          domain.CreateReviseItemDTO
		expectedItem domain.ReviseItem
	}{
		{
			name: "all fields filled",
			dto: domain.CreateReviseItemDTO{
				UserID:      gofakeit.UUID(),
				Name:        gofakeit.BookTitle(),
				Description: gofakeit.Sentence(10),
				Tags:        []string{gofakeit.HipsterWord(), gofakeit.HipsterWord()},
			},
		},
		{
			name: "no tags",
			dto: domain.CreateReviseItemDTO{
				UserID:      gofakeit.UUID(),
				Name:        gofakeit.BookTitle(),
				Description: gofakeit.Sentence(10),
			},
		},
		{
			name: "no description",
			dto: domain.CreateReviseItemDTO{
				UserID: gofakeit.UUID(),
				Name:   gofakeit.BookTitle(),
				Tags:   []string{gofakeit.HipsterWord(), gofakeit.HipsterWord()},
			},
		},
		{
			name: "no tags and description",
			dto: domain.CreateReviseItemDTO{
				UserID: gofakeit.UUID(),
				Name:   gofakeit.BookTitle(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := s.Service.Create(ctx, tt.dto)

			s.LogHandler.AssertEmpty()

			require.NoError(t, err)

			assert.NotEmpty(t, item.ID)
			assert.Equal(t, tt.dto.UserID, item.UserID.String())
			assert.Equal(t, tt.dto.Name, item.Name)
			assert.Equal(t, tt.dto.Description, item.Description)
			assert.Equal(t, tt.dto.Tags, []string(item.Tags))
			assert.Equal(t, domain.ReviseIteration(0), item.Iteration)

			getItem, err := s.Service.Get(ctx, item.ID.String())

			require.NoError(t, err)

			assert.Equal(t, item.ID.String(), getItem.ID.String())
			assert.Equal(t, item.UserID.String(), getItem.UserID.String())
			assert.Equal(t, item.Name, getItem.Name)
			assert.Equal(t, item.Description, getItem.Description)
			assert.Equal(t, item.Tags, getItem.Tags)
			assert.Equal(t, item.Iteration, getItem.Iteration)

		})
	}

}

func TestService_Revise_Create_FailPath(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name          string
		dto           domain.CreateReviseItemDTO
		expectedError error
	}{
		{
			name:          "empty user id",
			dto:           domain.CreateReviseItemDTO{},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "empty name",
			dto:           domain.CreateReviseItemDTO{UserID: gofakeit.UUID()},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "invalid user id",
			dto:           domain.CreateReviseItemDTO{UserID: "invalid", Name: gofakeit.BookTitle()},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "invalid name",
			dto:           domain.CreateReviseItemDTO{UserID: gofakeit.UUID(), Name: ""},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "invalid tags",
			dto:           domain.CreateReviseItemDTO{UserID: gofakeit.UUID(), Name: gofakeit.BookTitle(), Tags: []string{"tag1", ""}},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "invalid description",
			dto:           domain.CreateReviseItemDTO{UserID: gofakeit.UUID(), Name: gofakeit.BookTitle(), Description: "1"},
			expectedError: service.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Service.Create(ctx, tt.dto)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestService_Revise_Delete(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name   string
		id     string
		userID string
	}{
		{
			name:   "existing item",
			id:     "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID: "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
		},
		{
			name:   "existing item 2",
			id:     "4e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID: "2e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviseItem, err := s.Service.Delete(ctx, tt.id, tt.userID)

			require.NoError(t, err)

			assert.NotEmpty(t, reviseItem)

			_, err = s.Service.Get(ctx, tt.id)

			assert.ErrorIs(t, err, service.ErrNotFound)
		})
	}
}

func TestService_Revise_Delete_FailPath(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name          string
		id            string
		userID        string
		expectedError error
		exists        bool
	}{
		{
			name:          "non-existing item",
			id:            "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8f",
			userID:        "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			expectedError: service.ErrNotFound,
			exists:        false,
		},
		{
			name:          "invalid id",
			id:            "invalid",
			userID:        "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			expectedError: service.ErrInvalidArgument,
			exists:        false,
		},
		{
			name:          "empty id",
			id:            "",
			userID:        "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			expectedError: service.ErrInvalidArgument,
			exists:        false,
		},
		{
			name:          "invalid user id",
			id:            "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID:        "invalid",
			expectedError: service.ErrInvalidArgument,
			exists:        true,
		},
		{
			name:          "empty user id",
			id:            "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID:        "",
			expectedError: service.ErrInvalidArgument,
			exists:        true,
		},
		{
			name:          "user id mismatch",
			id:            "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID:        "2e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			expectedError: service.ErrUnauthorized,
			exists:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Service.Delete(ctx, tt.id, tt.userID)

			require.ErrorIs(t, err, tt.expectedError)

			if tt.exists {
				item, err := s.Service.Get(ctx, tt.id)

				require.NoError(t, err)
				assert.NotEmpty(t, item)
			}
		})
	}
}

func TestService_Revise_Update_HappyPath(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name         string
		id           string
		userID       string
		dto          domain.UpdateReviseItemDTO
		expectedItem domain.ReviseItem
	}{
		{
			name:   "all fields filled",
			id:     "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID: "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			dto: domain.UpdateReviseItemDTO{
				ID:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				UserID:       "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				Name:         gofakeit.BookTitle() + " ",                                     // " " for trimming test
				Description:  gofakeit.Sentence(10) + " ",                                    // " " for trimming test
				Tags:         []string{gofakeit.HipsterWord(), gofakeit.HipsterWord() + " "}, // " " for trimming test
				UpdateFields: []string{"name", "description", "tags"},
			},
		},
		{
			name: "name only",
			id:   "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			dto: domain.UpdateReviseItemDTO{
				ID:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				UserID:       "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				Name:         gofakeit.BookTitle(),
				UpdateFields: []string{"name"},
			},
		},
		{
			name: "description only",
			id:   "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			dto: domain.UpdateReviseItemDTO{
				ID:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				UserID:       "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				Description:  gofakeit.Sentence(10),
				UpdateFields: []string{"description"},
			},
		},
		{
			name: "tags only",
			id:   "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			dto: domain.UpdateReviseItemDTO{
				ID:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				UserID:       "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				Tags:         []string{gofakeit.HipsterWord(), gofakeit.HipsterWord()},
				UpdateFields: []string{"tags"},
			},
		},
		{
			name:   "description and tags",
			id:     "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID: "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			dto: domain.UpdateReviseItemDTO{
				ID:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				UserID:       "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				Description:  gofakeit.Sentence(10),
				Tags:         []string{gofakeit.HipsterWord(), gofakeit.HipsterWord()},
				UpdateFields: []string{"description", "tags"},
			},
		},
		{
			name:   "name and description",
			id:     "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID: "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			dto: domain.UpdateReviseItemDTO{
				ID:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				UserID:       "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				Name:         gofakeit.BookTitle(),
				Description:  gofakeit.Sentence(10),
				UpdateFields: []string{"name", "description"},
			},
		},
		{
			name:   "name and tags",
			id:     "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			userID: "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			dto: domain.UpdateReviseItemDTO{
				ID:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				UserID:       "1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
				Name:         gofakeit.BookTitle(),
				Tags:         []string{gofakeit.HipsterWord(), gofakeit.HipsterWord()},
				UpdateFields: []string{"name", "tags"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialItem, err := s.Service.Get(ctx, tt.id)
			require.NoError(t, err)

			item, err := s.Service.Update(ctx, tt.dto)

			s.LogHandler.AssertEmpty()

			require.NoError(t, err)

			trimmedTags := make([]string, 0, len(tt.dto.Tags))
			for _, tag := range tt.dto.Tags {
				trimmedTags = append(trimmedTags, strings.TrimSpace(tag))
			}

			assert.NotEmpty(t, item.ID)
			assert.NotEqual(t, initialItem.UpdatedAt, item.UpdatedAt)
			for _, field := range tt.dto.UpdateFields {
				switch field {
				case "name":
					assert.Equal(t, strings.TrimSpace(tt.dto.Name), item.Name)
				case "description":
					assert.Equal(t, strings.TrimSpace(tt.dto.Description), item.Description)
				case "tags":
					assert.Equal(t, trimmedTags, []string(item.Tags))
				}
			}

			getItem, err := s.Service.Get(ctx, item.ID.String())

			require.NoError(t, err)

			assert.Equal(t, item.ID.String(), getItem.ID.String())
			assert.Equal(t, item.UserID.String(), getItem.UserID.String())

			for _, field := range tt.dto.UpdateFields {
				switch field {
				case "name":
					assert.Equal(t, item.Name, getItem.Name)
				case "description":
					assert.Equal(t, item.Description, getItem.Description)
				case "tags":
					assert.Equal(t, item.Tags, getItem.Tags)
				}
			}

			if tt.dto.UpdateFields != nil || len(tt.dto.UpdateFields) > 0 {
				assert.NotEqual(t, item.UpdatedAt, getItem.UpdatedAt)
			} else {
				assert.Equal(t, item.UpdatedAt, getItem.UpdatedAt)
			}

			assert.Equal(t, item.Iteration, getItem.Iteration)
			assert.Equal(t, item.CreatedAt, getItem.CreatedAt)
			assert.Equal(t, item.LastRevisedAt, getItem.LastRevisedAt)
			assert.Equal(t, item.NextRevisionAt, getItem.NextRevisionAt)
		})
	}
}

func TestService_Revise_Update_FailPaht(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name          string
		dto           domain.UpdateReviseItemDTO
		expectedError error
	}{
		{
			name:          "empty id",
			dto:           domain.UpdateReviseItemDTO{},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "invalid id",
			dto:           domain.UpdateReviseItemDTO{ID: "invalid"},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "empty user id",
			dto:           domain.UpdateReviseItemDTO{ID: gofakeit.UUID()},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "invalid user id",
			dto:           domain.UpdateReviseItemDTO{ID: gofakeit.UUID(), UserID: "invalid"},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name:          "empty name",
			dto:           domain.UpdateReviseItemDTO{ID: gofakeit.UUID(), UserID: gofakeit.UUID()},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name: "empty description",
			dto: domain.UpdateReviseItemDTO{
				ID:     gofakeit.UUID(),
				UserID: gofakeit.UUID(),
				Name:   gofakeit.BookTitle(),
			},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name: "invalid tags",
			dto: domain.UpdateReviseItemDTO{
				ID:     gofakeit.UUID(),
				UserID: gofakeit.UUID(),
				Name:   gofakeit.BookTitle(),
				Tags:   []string{"tag1", ""},
			},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name: "invalid description length",
			dto: domain.UpdateReviseItemDTO{
				ID:          gofakeit.UUID(),
				UserID:      gofakeit.UUID(),
				Name:        gofakeit.BookTitle(),
				Description: "1",
			},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name: "invalid update fields",
			dto: domain.UpdateReviseItemDTO{
				ID:           gofakeit.UUID(),
				UserID:       gofakeit.UUID(),
				Name:         gofakeit.BookTitle(),
				Description:  gofakeit.Sentence(10),
				UpdateFields: []string{"invalid"},
			},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name: "empty update fields",
			dto: domain.UpdateReviseItemDTO{
				ID:           gofakeit.UUID(),
				UserID:       gofakeit.UUID(),
				Name:         gofakeit.BookTitle(),
				Description:  gofakeit.Sentence(10),
				UpdateFields: []string{},
			},
			expectedError: service.ErrInvalidArgument,
		},
		{
			name: "non-existing item",
			dto: domain.UpdateReviseItemDTO{
				ID:           gofakeit.UUID(),
				UserID:       gofakeit.UUID(),
				Name:         gofakeit.BookTitle(),
				Description:  gofakeit.Sentence(10),
				UpdateFields: []string{"name", "description"},
			},
			expectedError: service.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Service.Update(ctx, tt.dto)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
