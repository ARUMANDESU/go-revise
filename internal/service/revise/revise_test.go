package revisesvc

import (
	"context"
	"errors"
	"testing"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/ARUMANDESU/go-revise/internal/service/revise/mocks"
	"github.com/ARUMANDESU/go-revise/internal/storage"
	"github.com/ARUMANDESU/go-revise/pkg/logger"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofrs/uuid"
	guuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type Suite struct {
	service            Revise
	mockReviseProvider *mocks.ReviseProvider
	mockReviseManager  *mocks.ReviseManager
	mockUserProvider   *mocks.UserProvider
}

func NewSuite(t *testing.T) Suite {
	t.Helper()

	mockReviseProvider := mocks.NewReviseProvider(t)
	mockReviseManager := mocks.NewReviseManager(t)
	mockUserProvider := mocks.NewUserProvider(t)

	return Suite{
		service: NewRevise(logger.Plug(),
			ReviseStorages{
				ReviseProvider: mockReviseProvider,
				ReviseManager:  mockReviseManager,
				UserProvider:   mockUserProvider,
			},
		),
		mockReviseProvider: mockReviseProvider,
		mockReviseManager:  mockReviseManager,
		mockUserProvider:   mockUserProvider,
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
			id:         guuid.New().String(),
			provideErr: nil,
			wantErr:    nil,
		},
		{
			name:       "error: empty ID",
			id:         "",
			provideErr: nil,
			wantErr:    service.ErrInvalidArgument,
		},
		{
			name:       "error: revise not found",
			id:         guuid.New().String(),
			provideErr: storage.ErrNotFound,
			wantErr:    service.ErrNotFound,
		},
		{
			name:       "error: internal error",
			id:         guuid.New().String(),
			provideErr: errors.New("unexpected db error"),
			wantErr:    service.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.wantErr, service.ErrInvalidArgument) {
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
				UserID:      guuid.New().String(),
				Name:        gofakeit.LetterN(domain.ValidNameMinLength + 1),
				Tags:        []string{gofakeit.LetterN(domain.ValidTagsMinLength + 1)},
				Description: gofakeit.Sentence(domain.ValidDescriptionMinLength + 1),
			},
			mockErr: nil,
			wantErr: nil,
		},
		{
			name: "success: empty tags and description",
			dto: domain.CreateReviseItemDTO{
				UserID: guuid.New().String(),
				Name:   gofakeit.LetterN(domain.ValidNameMinLength + 1),
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
			wantErr: service.ErrInvalidArgument,
		},
		{
			name: "error: internal error",
			dto: domain.CreateReviseItemDTO{
				UserID: guuid.New().String(),
				Name:   gofakeit.LetterN(domain.ValidNameMinLength + 1),
			},
			mockErr: errors.New("unexpected db error"),
			wantErr: service.ErrInternal,
		},
		{
			name: "error: invalid arguments",
			dto: domain.CreateReviseItemDTO{
				UserID:      "",
				Name:        gofakeit.LetterN(domain.ValidNameMinLength - 1),
				Tags:        []string{gofakeit.LetterN(domain.ValidTagsMinLength - 1)},
				Description: gofakeit.Sentence(domain.ValidDescriptionMinLength - 1),
			},
			mockErr: nil,
			wantErr: service.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.wantErr, service.ErrInvalidArgument) {
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

func TestRevise_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := NewSuite(t)
		revisionID := guuid.New().String()
		userID := guuid.New().String()
		expectedItem := domain.ReviseItem{
			ID:     uuid.FromStringOrNil(revisionID),
			UserID: uuid.FromStringOrNil(userID),
			Name:   gofakeit.LetterN(domain.ValidNameMinLength + 1),
		}

		s.mockReviseProvider.On("GetRevise", mock.Anything, revisionID).Return(expectedItem, nil)
		defer s.mockReviseProvider.AssertExpectations(t)
		s.mockReviseManager.On("DeleteRevise", mock.Anything, revisionID).Return(nil)
		defer s.mockReviseManager.AssertExpectations(t)

		gotItem, err := s.service.Delete(context.Background(), revisionID, userID)

		require.NoError(t, err)

		assert.Equal(t, expectedItem, gotItem)
	})
}

func TestRevise_Delete_FailPath(t *testing.T) {

	revisionID := guuid.New().String()
	userID := guuid.New().String()

	tests := []struct {
		name        string
		revisionID  string
		userID      string
		reviseItem  domain.ReviseItem
		onGetErr    error
		onDeleteErr error
		wantErr     error
	}{
		{
			name:        "error: empty ID",
			revisionID:  "",
			userID:      userID,
			reviseItem:  domain.ReviseItem{},
			onGetErr:    nil,
			onDeleteErr: nil,
			wantErr:     service.ErrInvalidArgument,
		},
		{
			name:        "error: empty user ID",
			revisionID:  revisionID,
			userID:      "",
			reviseItem:  domain.ReviseItem{},
			onGetErr:    nil,
			onDeleteErr: nil,
			wantErr:     service.ErrInvalidArgument,
		},
		{
			name:        "error: revise not found",
			revisionID:  revisionID,
			userID:      userID,
			reviseItem:  domain.ReviseItem{},
			onGetErr:    storage.ErrNotFound,
			onDeleteErr: nil,
			wantErr:     service.ErrNotFound,
		},
		{
			name:       "error: not found on delete",
			revisionID: revisionID,
			userID:     userID,
			reviseItem: domain.ReviseItem{
				ID:     uuid.FromStringOrNil(revisionID),
				UserID: uuid.FromStringOrNil(userID),
			},
			onGetErr:    nil,
			onDeleteErr: storage.ErrNotFound,
			wantErr:     service.ErrNotFound,
		},
		{
			name:       "error: unauthorized",
			revisionID: revisionID,
			userID:     userID,
			reviseItem: domain.ReviseItem{
				ID:     uuid.FromStringOrNil(revisionID),
				UserID: uuid.FromStringOrNil(guuid.New().String()), // different user ID
			},
			onGetErr:    nil,
			onDeleteErr: nil,
			wantErr:     service.ErrUnauthorized,
		},
		{
			name:        "error: internal error",
			revisionID:  revisionID,
			userID:      userID,
			reviseItem:  domain.ReviseItem{},
			onGetErr:    errors.New("unexpected db error"),
			onDeleteErr: nil,
			wantErr:     service.ErrInternal,
		},
		{
			name:       "error: internal error",
			revisionID: revisionID,
			userID:     userID,
			// this is not empty to not get unauthorized error
			// because first "get" is successful and then it checks the user ID
			// and then it tries to delete it and fails
			reviseItem: domain.ReviseItem{
				ID:     uuid.FromStringOrNil(revisionID),
				UserID: uuid.FromStringOrNil(userID),
			},
			onGetErr:    nil,
			onDeleteErr: errors.New("unexpected db error"),
			wantErr:     service.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.wantErr, service.ErrInvalidArgument) {
				s.mockReviseProvider.On("GetRevise", mock.Anything, tt.revisionID).Return(tt.reviseItem, tt.onGetErr)
				defer s.mockReviseProvider.AssertExpectations(t)
				if tt.onGetErr == nil && !errors.Is(tt.wantErr, service.ErrUnauthorized) {
					s.mockReviseManager.On("DeleteRevise", mock.Anything, tt.revisionID).Return(tt.onDeleteErr)
					defer s.mockReviseManager.AssertExpectations(t)
				}
			}

			_, err := s.service.Delete(context.Background(), tt.revisionID, tt.userID)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}

}

func TestRevise_Update_HappyPath(t *testing.T) {
	uid, _ := uuid.NewV7()
	revisionID, _ := uuid.NewV7()

	tests := []struct {
		name    string
		dto     domain.UpdateReviseItemDTO
		initial domain.ReviseItem
	}{
		{
			name: "success",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				Tags:         []string{gofakeit.LetterN(domain.ValidTagsMinLength + 1)},
				Description:  gofakeit.Sentence(domain.ValidDescriptionMinLength + 1),
				UpdateFields: []string{"name", "description", "tags"},
			},
			initial: domain.ReviseItem{
				ID:          revisionID,
				UserID:      uid,
				Name:        gofakeit.LetterN(domain.ValidNameMinLength + 1),
				Tags:        domain.StringArray{gofakeit.LetterN(domain.ValidTagsMinLength + 1)},
				Description: gofakeit.Sentence(domain.ValidDescriptionMinLength + 1),
			},
		},
		{
			name: "success: update name",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uid,
				Name:   gofakeit.LetterN(domain.ValidNameMinLength + 1),
			},
		},
		{
			name: "success: update description",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Description:  gofakeit.Sentence(domain.ValidDescriptionMinLength + 1),
				UpdateFields: []string{"description"},
			},
			initial: domain.ReviseItem{
				ID:          revisionID,
				UserID:      uid,
				Description: gofakeit.Sentence(domain.ValidDescriptionMinLength + 1),
			},
		},
		{
			name: "success: update set description to empty",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				UpdateFields: []string{"description"},
			},
			initial: domain.ReviseItem{
				ID:          revisionID,
				UserID:      uid,
				Description: gofakeit.Sentence(domain.ValidDescriptionMinLength + 1),
			},
		},
		{
			name: "success: update tags",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Tags:         []string{gofakeit.LetterN(domain.ValidTagsMinLength + 1)},
				UpdateFields: []string{"tags"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uid,
				Tags:   domain.StringArray{gofakeit.LetterN(domain.ValidTagsMinLength + 1)},
			},
		},
		{
			name: "success: update set tags to empty",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				UpdateFields: []string{"tags"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uid,
				Tags:   domain.StringArray{gofakeit.LetterN(domain.ValidTagsMinLength + 1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			s.mockReviseProvider.On("GetRevise", mock.Anything, tt.dto.ID).Return(tt.initial, nil)
			defer s.mockReviseProvider.AssertExpectations(t)

			s.mockReviseManager.On("UpdateRevise", mock.Anything, mock.AnythingOfType("domain.ReviseItem")).Return(nil)
			defer s.mockReviseManager.AssertExpectations(t)

			reviseItem, err := s.service.Update(context.Background(), tt.dto)

			require.NoError(t, err)

			assert.Equal(t, tt.dto.ID, reviseItem.ID.String())
			assert.Equal(t, tt.dto.UserID, reviseItem.UserID.String())
			for _, field := range tt.dto.UpdateFields {
				switch field {
				case "name":
					assert.Equal(t, tt.dto.Name, reviseItem.Name)
				case "tags":
					assert.Equal(t, domain.StringArray(tt.dto.Tags), reviseItem.Tags)
				case "description":
					assert.Equal(t, tt.dto.Description, reviseItem.Description)
				}
			}
		})
	}
}

func TestRevise_Update_FailPath(t *testing.T) {
	uid, _ := uuid.NewV7()
	revisionID, _ := uuid.NewV7()

	tests := []struct {
		name        string
		dto         domain.UpdateReviseItemDTO
		initial     domain.ReviseItem
		onGetErr    error
		onUpdateErr error
		wantErr     error
	}{
		{
			name: "error: empty ID",
			dto: domain.UpdateReviseItemDTO{
				ID:           "",
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial:     domain.ReviseItem{},
			onGetErr:    nil,
			onUpdateErr: nil,
			wantErr:     service.ErrInvalidArgument,
		},
		{
			name: "error: empty user ID",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       "",
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial:     domain.ReviseItem{},
			onGetErr:    nil,
			onUpdateErr: nil,
			wantErr:     service.ErrInvalidArgument,
		},
		{
			name: "error: revise not found",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial:     domain.ReviseItem{},
			onGetErr:    storage.ErrNotFound,
			onUpdateErr: nil,
			wantErr:     service.ErrNotFound,
		},
		{
			name: "error: not found on update",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uid,
			},
			onGetErr:    nil,
			onUpdateErr: storage.ErrNotFound,
			wantErr:     service.ErrNotFound,
		},
		{
			name: "error: unauthorized",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uuid.FromStringOrNil(guuid.New().String()), // different user ID
			},
			onGetErr:    nil,
			wantErr:     service.ErrUnauthorized,
			onUpdateErr: nil,
		},
		{
			name: "error: internal error",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial:     domain.ReviseItem{},
			onGetErr:    errors.New("unexpected db error"),
			onUpdateErr: nil,
			wantErr:     service.ErrInternal,
		},
		{
			name: "error: internal error",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"name"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uid,
				Name:   gofakeit.LetterN(domain.ValidNameMinLength + 1),
			},
			onGetErr:    nil,
			onUpdateErr: errors.New("unexpected db error"),
			wantErr:     service.ErrInternal,
		},
		{
			name: "error: invalid arguments",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength - 1),
				UpdateFields: []string{"name"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uid,
				Name:   gofakeit.LetterN(domain.ValidNameMinLength + 1),
			},
			onGetErr:    nil,
			onUpdateErr: nil,
			wantErr:     service.ErrInvalidArgument,
		},
		{
			name: "error: invalid update fields",
			dto: domain.UpdateReviseItemDTO{
				ID:           revisionID.String(),
				UserID:       uid.String(),
				Name:         gofakeit.LetterN(domain.ValidNameMinLength + 1),
				UpdateFields: []string{"invalid"},
			},
			initial: domain.ReviseItem{
				ID:     revisionID,
				UserID: uid,
				Name:   gofakeit.LetterN(domain.ValidNameMinLength + 1),
			},
			onGetErr:    nil,
			onUpdateErr: nil,
			wantErr:     service.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.wantErr, service.ErrInvalidArgument) {
				s.mockReviseProvider.On("GetRevise", mock.Anything, tt.dto.ID).Return(tt.initial, tt.onGetErr)
				defer s.mockReviseProvider.AssertExpectations(t)
				if tt.onGetErr == nil && !errors.Is(tt.wantErr, service.ErrUnauthorized) {
					s.mockReviseManager.On("UpdateRevise", mock.Anything, mock.AnythingOfType("domain.ReviseItem")).Return(tt.onUpdateErr)
					defer s.mockReviseManager.AssertExpectations(t)
				}
			}

			_, err := s.service.Update(context.Background(), tt.dto)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRevise_List_HappyPath(t *testing.T) {
	uid, _ := uuid.NewV7()

	tests := []struct {
		name string
		dto  domain.ListReviseItemDTO
	}{
		{
			name: "success",
			dto: domain.ListReviseItemDTO{
				UserID: uid.String(),
				Pagination: &domain.Pagination{
					Page:     1,
					PageSize: 10,
				},
				Sort: &domain.Sort{
					Field: domain.SortFieldDefault,
					Order: domain.SortOrderAsc,
				},
			},
		},
		{
			name: "success: telegram user ID",
			dto: domain.ListReviseItemDTO{
				UserID: int64(123),
				Pagination: &domain.Pagination{
					Page:     1,
					PageSize: 10,
				},
				Sort: &domain.Sort{
					Field: domain.SortFieldDefault,
					Order: domain.SortOrderAsc,
				},
			},
		},
		{
			name: "success: no pagination",
			dto: domain.ListReviseItemDTO{
				UserID: uid.String(),
				Sort: &domain.Sort{
					Field: domain.SortFieldDefault,
					Order: domain.SortOrderAsc,
				},
			},
		},
		{
			name: "success: no sort",
			dto: domain.ListReviseItemDTO{
				UserID: uid.String(),
				Pagination: &domain.Pagination{
					Page:     1,
					PageSize: 10,
				},
			},
		},
		{
			name: "success: no pagination and sort",
			dto: domain.ListReviseItemDTO{
				UserID: uid.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if _, ok := tt.dto.UserID.(int64); ok {
				s.mockUserProvider.On("GetUserByTelegramID", mock.Anything, tt.dto.UserID.(int64)).
					Return(domain.User{ID: uid}, nil)
				defer s.mockUserProvider.AssertExpectations(t)
			}
			s.mockReviseProvider.On("ListRevises", mock.Anything, mock.AnythingOfType("domain.ListReviseItemDTO")).
				Return([]domain.ReviseItem{}, domain.PaginationMetadata{}, nil)
			defer s.mockReviseProvider.AssertExpectations(t)

			_, _, err := s.service.List(context.Background(), tt.dto)

			require.NoError(t, err, "got %v", err)
		})
	}
}

func TestRevise_List_FailPath(t *testing.T) {
	uid, _ := uuid.NewV7()

	tests := []struct {
		name      string
		dto       domain.ListReviseItemDTO
		onGetErr  error
		onListErr error
		wantErr   error
	}{
		{
			name: "error: empty user ID",
			dto: domain.ListReviseItemDTO{
				UserID: "",
			},
			onGetErr:  nil,
			onListErr: nil,
			wantErr:   service.ErrInvalidArgument,
		},
		{
			name: "error: invalid user ID",
			dto: domain.ListReviseItemDTO{
				UserID: "invalid",
			},
			onGetErr:  nil,
			onListErr: nil,
			wantErr:   service.ErrInvalidArgument,
		},
		{
			name: "error: user not found",
			dto: domain.ListReviseItemDTO{
				UserID: int64(123),
			},
			onGetErr:  storage.ErrNotFound,
			onListErr: nil,
			wantErr:   service.ErrNotFound,
		},
		{
			name: "error: revise not found",
			dto: domain.ListReviseItemDTO{
				UserID: uid.String(),
			},
			onGetErr:  nil,
			onListErr: storage.ErrNotFound,
			wantErr:   service.ErrNotFound,
		},
		{
			name: "error: internal error",
			dto: domain.ListReviseItemDTO{
				UserID: int64(123),
			},
			onGetErr:  errors.New("unexpected db error"),
			onListErr: nil,
			wantErr:   service.ErrInternal,
		},
		{
			name: "error: internal error",
			dto: domain.ListReviseItemDTO{
				UserID: uid.String(),
			},
			onGetErr:  nil,
			onListErr: errors.New("unexpected db error"),
			wantErr:   service.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.wantErr, service.ErrInvalidArgument) {
				if _, ok := tt.dto.UserID.(int64); ok {
					s.mockUserProvider.On("GetUserByTelegramID", mock.Anything, mock.AnythingOfType("int64")).
						Return(domain.User{}, tt.onGetErr)
					defer s.mockUserProvider.AssertExpectations(t)
				}
				if tt.onGetErr == nil {
					s.mockReviseProvider.On("ListRevises", mock.Anything, mock.AnythingOfType("domain.ListReviseItemDTO")).
						Return([]domain.ReviseItem{}, domain.PaginationMetadata{}, tt.onListErr)
					defer s.mockReviseProvider.AssertExpectations(t)
				}
			}

			_, _, err := s.service.List(context.Background(), tt.dto)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
