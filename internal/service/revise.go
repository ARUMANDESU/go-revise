package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/storage"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid"
)

//go:generate mockery --name ReviseProvider --output mocks
type ReviseProvider interface {
	GetRevise(ctx context.Context, id string) (domain.ReviseItem, error)
	ListRevises(ctx context.Context, userID string) ([]domain.ReviseItem, domain.PaginationMetadata, error)
}

//go:generate mockery --name ReviseManager --output mocks
type ReviseManager interface {
	CreateRevise(ctx context.Context, revise domain.ReviseItem) error
	UpdateRevise(ctx context.Context, revise domain.ReviseItem) (domain.ReviseItem, error)
	DeleteRevise(ctx context.Context, id string) (domain.ReviseItem, error)
}

type Revise struct {
	log            *slog.Logger
	reviseProvider ReviseProvider
	reviseManager  ReviseManager
}

func NewRevise(log *slog.Logger, reviseProvider ReviseProvider, reviseManager ReviseManager) Revise {
	return Revise{
		log:            log,
		reviseProvider: reviseProvider,
		reviseManager:  reviseManager,
	}
}

func (r *Revise) Get(ctx context.Context, id string) (domain.ReviseItem, error) {
	const op = "service.revise.get"

	err := validation.Validate(id, validation.Required)
	if err != nil {
		return domain.ReviseItem{}, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	revise, err := r.reviseProvider.GetRevise(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.ReviseItem{}, ErrNotFound
		default:
			r.log.Error(domain.WrapErrorWithOp(err, op, "failed to get revise").Error())
			return domain.ReviseItem{}, ErrInternal
		}
	}

	return revise, nil
}

func (r *Revise) List(ctx context.Context, userID string) ([]domain.ReviseItem, domain.PaginationMetadata, error) {
	const op = "service.revise.list"
	panic("not implemented") // TODO: Implement
}

func (r *Revise) Create(ctx context.Context, dto domain.CreateReviseItemDTO) (domain.ReviseItem, error) {
	const op = "service.revise.create"

	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.UserID, validation.Required),
		validation.Field(&dto.Name, validation.Required, validation.By(validateName)),
		validation.Field(&dto.Tags, validation.By(validateTags)),
		validation.Field(&dto.Description, validation.By(validateDescription)),
	)
	if err != nil {
		return domain.ReviseItem{}, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		r.log.Error(domain.WrapErrorWithOp(err, op, "failed to generate new UUID").Error())
		return domain.ReviseItem{}, ErrInternal
	}

	revise := domain.ReviseItem{
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

	if err := r.reviseManager.CreateRevise(ctx, revise); err != nil {
		r.log.Error(domain.WrapErrorWithOp(err, op, "failed to create revise").Error())
		return domain.ReviseItem{}, ErrInternal
	}

	return revise, nil
}

func (r *Revise) Update(ctx context.Context, revise domain.ReviseItem) (domain.ReviseItem, error) {
	const op = "service.revise.update"
	panic("not implemented") // TODO: Implement
}

func (r *Revise) Delete(ctx context.Context, id string) (domain.ReviseItem, error) {
	const op = "service.revise.delete"
	panic("not implemented") // TODO: Implement
}
