package revisesvc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/ARUMANDESU/go-revise/internal/storage"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid"
)

//go:generate mockery --name ReviseProvider --output mocks
type ReviseProvider interface {
	GetRevise(ctx context.Context, id string) (domain.ReviseItem, error)
	ListRevises(ctx context.Context, dto domain.ListReviseItemDTO) ([]domain.ReviseItem, domain.PaginationMetadata, error)
}

//go:generate mockery --name ReviseManager --output mocks
type ReviseManager interface {
	CreateRevise(ctx context.Context, revise domain.ReviseItem) error
	UpdateRevise(ctx context.Context, revise domain.ReviseItem) error
	DeleteRevise(ctx context.Context, id string) error
}

//go:generate mockery --name UserProvider --output mocks
type UserProvider interface {
	GetUser(ctx context.Context, id string) (domain.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (domain.User, error)
}

type ReviseStorages struct {
	ReviseProvider ReviseProvider
	ReviseManager  ReviseManager
	UserProvider   UserProvider
}

type Revise struct {
	log *slog.Logger
	ReviseStorages
}

func NewRevise(log *slog.Logger, storages ReviseStorages) Revise {
	return Revise{
		log:            log,
		ReviseStorages: storages,
	}
}

func (r *Revise) Get(ctx context.Context, id string) (domain.ReviseItem, error) {
	const op = "service.revise.get"

	err := validation.Validate(id, validation.Required, is.UUID)
	if err != nil {
		return domain.ReviseItem{}, fmt.Errorf("%w: %w", service.ErrInvalidArgument, err)
	}

	reviseItem, err := r.ReviseProvider.GetRevise(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.ReviseItem{}, service.ErrNotFound
		default:
			r.log.Error(domain.WrapErrorWithOp(err, op, "failed to get revise").Error())
			return domain.ReviseItem{}, service.ErrInternal
		}
	}

	return reviseItem, nil
}

func (r *Revise) List(ctx context.Context, dto domain.ListReviseItemDTO) ([]domain.ReviseItem, domain.PaginationMetadata, error) {
	const op = "service.revise.list"

	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.UserID, validation.Required, validation.By(domain.ValidateFilterUserID)),
		validation.Field(&dto.Pagination, validation.By(domain.ValidatePagination)),
		validation.Field(&dto.Sort, validation.By(domain.ValidateSort)),
	)
	if err != nil {
		return nil, domain.PaginationMetadata{}, fmt.Errorf("%w: %w", service.ErrInvalidArgument, err)
	}

	// if user id is int64, then get the string(uuid) from the database
	if telegramID, ok := dto.UserID.(int64); ok {
		user, err := r.UserProvider.GetUserByTelegramID(ctx, telegramID)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrNotFound):
				return nil, domain.PaginationMetadata{}, service.ErrNotFound
			default:
				r.log.Error(domain.WrapErrorWithOp(err, op, "failed to get user").Error())
				return nil, domain.PaginationMetadata{}, service.ErrInternal
			}
		}
		dto.UserID = user.ID.String()
	}

	if dto.Pagination == nil {
		dto.Pagination = domain.DefaultPagination()
	}
	if dto.Sort == nil {
		dto.Sort = domain.DefaultSort()
	}

	reviseItems, pagination, err := r.ReviseProvider.ListRevises(ctx, dto)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, domain.PaginationMetadata{}, service.ErrNotFound
		}
		r.log.Error(domain.WrapErrorWithOp(err, op, "failed to list revise").Error())
		return nil, domain.PaginationMetadata{}, service.ErrInternal
	}

	return reviseItems, pagination, nil
}

func (r *Revise) Create(ctx context.Context, dto domain.CreateReviseItemDTO) (domain.ReviseItem, error) {
	const op = "service.revise.create"

	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.UserID, validation.Required, is.UUID),
		validation.Field(&dto.Name, validation.Required, validation.By(domain.ValidateName)),
		validation.Field(&dto.Tags, validation.By(domain.ValidateTags)),
		validation.Field(&dto.Description, validation.By(domain.ValidateDescription)),
	)
	if err != nil {
		return domain.ReviseItem{}, fmt.Errorf("%w: %w", service.ErrInvalidArgument, err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		r.log.Error(domain.WrapErrorWithOp(err, op, "failed to generate new UUID").Error())
		return domain.ReviseItem{}, service.ErrInternal
	}

	reviseItem := domain.ReviseItem{
		ID:             id,
		UserID:         uuid.FromStringOrNil(dto.UserID),
		Name:           dto.Name,
		Tags:           dto.Tags,
		Description:    dto.Description,
		Iteration:      0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastRevisedAt:  time.Now(),
		NextRevisionAt: time.Now().Add(time.Duration(domain.IntervalMap[0])),
	}

	// TODO: create reminder

	if err := r.ReviseManager.CreateRevise(ctx, reviseItem); err != nil {
		r.log.Error(domain.WrapErrorWithOp(err, op, "failed to create revise").Error())
		return domain.ReviseItem{}, service.ErrInternal
	}

	return reviseItem, nil
}

func (r *Revise) Update(ctx context.Context, dto domain.UpdateReviseItemDTO) (domain.ReviseItem, error) {
	const op = "service.revise.update"

	reviseItemUpdateFields := []any{"name", "description", "tags"}

	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, validation.Required, is.UUID),
		validation.Field(&dto.UserID, validation.Required, is.UUID),
		validation.Field(&dto.Name, validation.By(domain.ValidateName)),
		validation.Field(&dto.Tags, validation.By(domain.ValidateTags)),
		validation.Field(&dto.Description, validation.By(domain.ValidateDescription)),
		validation.Field(&dto.UpdateFields, validation.Required, validation.Each(validation.In(reviseItemUpdateFields...))),
	)
	if err != nil {
		return domain.ReviseItem{}, fmt.Errorf("%w: %w", service.ErrInvalidArgument, err)
	}

	reviseItem, err := r.ReviseProvider.GetRevise(ctx, dto.ID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.ReviseItem{}, service.ErrNotFound
		default:
			r.log.Error(domain.WrapErrorWithOp(err, op, "failed to get revise").Error())
			return domain.ReviseItem{}, service.ErrInternal
		}
	}

	if !reviseItem.AbleToUpdate(dto.UserID) {
		return domain.ReviseItem{}, service.ErrUnauthorized
	}

	reviseItem = reviseItem.PartialUpdate(dto)

	if err := r.ReviseManager.UpdateRevise(ctx, reviseItem); err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.ReviseItem{}, service.ErrNotFound
		default:
			r.log.Error(domain.WrapErrorWithOp(err, op, "failed to update revise").Error())
			return domain.ReviseItem{}, service.ErrInternal
		}
	}

	return reviseItem, nil
}

func (r *Revise) Delete(ctx context.Context, id string, userID string) (domain.ReviseItem, error) {
	const op = "service.revise.delete"

	err := validation.Validate(id, validation.Required.Error("id must provided"), is.UUID.Error("id must be a valid UUID"))
	if err != nil {
		return domain.ReviseItem{}, fmt.Errorf("%w: %w", service.ErrInvalidArgument, err)
	}
	err = validation.Validate(userID, validation.Required.Error("userID must provided"), is.UUID.Error("userID must be a valid UUID"))
	if err != nil {
		return domain.ReviseItem{}, fmt.Errorf("%w: %w", service.ErrInvalidArgument, err)
	}

	reviseItem, err := r.ReviseProvider.GetRevise(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.ReviseItem{}, service.ErrNotFound
		default:
			r.log.Error(domain.WrapErrorWithOp(err, op, "failed to get revise").Error())
			return domain.ReviseItem{}, service.ErrInternal
		}
	}

	if !reviseItem.AbleToUpdate(userID) {
		return domain.ReviseItem{}, service.ErrUnauthorized
	}

	err = r.ReviseManager.DeleteRevise(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.ReviseItem{}, service.ErrNotFound
		default:
			r.log.Error(domain.WrapErrorWithOp(err, op, "failed to get revise").Error())
			return domain.ReviseItem{}, service.ErrInternal
		}
	}

	return reviseItem, nil
}
