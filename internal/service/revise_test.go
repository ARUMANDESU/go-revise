package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service/mocks"
	"github.com/ARUMANDESU/go-revise/internal/storage"
	"github.com/ARUMANDESU/go-revise/pkg/logger"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type Suite struct {
	service            Revise
	mockReviseProvider *mocks.ReviseProvider
	mockReviseManager  *mocks.ReviseManager
}

func NewSuite(t *testing.T) Suite {
	t.Helper()

	mockReviseProvider := mocks.NewReviseProvider(t)
	mockReviseManager := mocks.NewReviseManager(t)

	return Suite{
		service:            NewRevise(logger.Plug(), mockReviseProvider, mockReviseManager),
		mockReviseProvider: mockReviseProvider,
		mockReviseManager:  mockReviseManager,
	}
}

func TestRevise_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		provideErr error
		wantErr    error
	}{
		{
			name:       "success",
			id:         uuid.New().String(),
			provideErr: nil,
			wantErr:    nil,
		},
		{
			name:       "error: empty ID",
			id:         "",
			provideErr: nil,
			wantErr:    ErrInvalidArgument,
		},
		{
			name:       "error: revise not found",
			id:         uuid.New().String(),
			provideErr: storage.ErrNotFound,
			wantErr:    ErrNotFound,
		},
		{
			name:       "error: internal error",
			id:         uuid.New().String(),
			provideErr: errors.New("unexpected db error"),
			wantErr:    ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.wantErr, ErrInvalidArgument) {
				s.mockReviseProvider.On("GetRevise", mock.Anything, tt.id).Return(domain.ReviseItem{}, tt.provideErr)
				defer s.mockReviseProvider.AssertExpectations(t)
			}

			_, err := s.service.Get(context.Background(), tt.id)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRevise_Create(t *testing.T) {
	tests := []struct {
		name    string
		dto     domain.CreateReviseItemDTO
		mockErr error
		wantErr error
	}{
		{
			name: "success",
			dto: domain.CreateReviseItemDTO{
				UserID:      uuid.New().String(),
				Name:        gofakeit.LetterN(ValidNameMinLength + 1),
				Tags:        []string{gofakeit.LetterN(ValidTagsMinLength + 1)},
				Description: gofakeit.Sentence(ValidDescriptionMinLength + 1),
			},
			mockErr: nil,
			wantErr: nil,
		},
		{
			name: "success: empty tags and description",
			dto: domain.CreateReviseItemDTO{
				UserID: uuid.New().String(),
				Name:   gofakeit.LetterN(ValidNameMinLength + 1),
			},
			mockErr: nil,
			wantErr: nil,
		},
		{
			name: "error: empty revise",
			dto: domain.CreateReviseItemDTO{
				UserID: "", // must be provided
				Name:   "", // must by non-empty
			},
			mockErr: nil,
			wantErr: ErrInvalidArgument,
		},
		{
			name: "error: internal error",
			dto: domain.CreateReviseItemDTO{
				UserID: uuid.New().String(),
				Name:   gofakeit.LetterN(ValidNameMinLength + 1),
			},
			mockErr: errors.New("unexpected db error"),
			wantErr: ErrInternal,
		},
		{
			name: "error: invalid arguments",
			dto: domain.CreateReviseItemDTO{
				UserID:      "",
				Name:        gofakeit.LetterN(ValidNameMinLength - 1),
				Tags:        []string{gofakeit.LetterN(ValidTagsMinLength - 1)},
				Description: gofakeit.Sentence(ValidDescriptionMinLength - 1),
			},
			mockErr: nil,
			wantErr: ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.wantErr, ErrInvalidArgument) {
				s.mockReviseManager.On("CreateRevise", mock.Anything, mock.AnythingOfType("domain.ReviseItem")).Return(tt.mockErr)
				defer s.mockReviseManager.AssertExpectations(t)
			}

			reviseItem, err := s.service.Create(context.Background(), tt.dto)

			require.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.Equal(t, tt.dto.UserID, reviseItem.UserID.String())
				assert.Equal(t, tt.dto.Name, reviseItem.Name)
				assert.Equal(t, domain.StringArray(tt.dto.Tags), reviseItem.Tags)
				assert.Equal(t, tt.dto.Description, reviseItem.Description)
				assert.Equal(t, domain.ReviseIteration(0), reviseItem.Iteration)
			}
		})
	}
}
