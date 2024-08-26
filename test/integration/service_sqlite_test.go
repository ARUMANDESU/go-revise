package integration

import (
	"context"
	"testing"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_Revise_Get(t *testing.T) {
	ctx := context.Background()
	s, cleanup := NewSuite(t)
	defer cleanup()

	tests := []struct {
		name         string
		id           string
		wantErr      error
		expectedItem domain.ReviseItem
	}{
		{
			name:    "get existing item",
			id:      "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e",
			wantErr: nil,
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
			name:         "get non-existing item",
			id:           "3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8f",
			wantErr:      service.ErrNotFound,
			expectedItem: domain.ReviseItem{},
		},
		{
			name:         "get invalid item",
			id:           "invalid",
			wantErr:      service.ErrInvalidArgument,
			expectedItem: domain.ReviseItem{},
		},
		{
			name:         "get empty item",
			id:           "",
			wantErr:      service.ErrInvalidArgument,
			expectedItem: domain.ReviseItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := s.Service.Get(ctx, tt.id)

			assert.ErrorIs(t, err, tt.wantErr)

			assert.Equal(t, tt.expectedItem.ID.String(), item.ID.String())
			assert.Equal(t, tt.expectedItem.UserID.String(), item.UserID.String())
			assert.Equal(t, tt.expectedItem.Name, item.Name)
			assert.Equal(t, tt.expectedItem.Description, item.Description)
			assert.Equal(t, tt.expectedItem.Tags, item.Tags)
			assert.Equal(t, tt.expectedItem.Iteration, item.Iteration)
		})
	}
}
