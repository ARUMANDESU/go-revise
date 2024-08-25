package service

import (
	"context"
	"log/slog"

	"github.com/ARUMANDESU/go-revise/internal/domain"
)

type ReviseProvider interface {
	GetRevise(ctx context.Context, id string) (domain.ReviseItem, error)
	ListRevises(ctx context.Context, userID string) ([]domain.ReviseItem, domain.PaginationMetadata, error)
}

type ReviseManager interface {
	CreateRevise(ctx context.Context, revise domain.ReviseItem) (string, error)
	UpdateRevise(ctx context.Context, revise domain.ReviseItem) error
	DeleteRevise(ctx context.Context, id string) error
}

type Revise struct {
	log            *slog.Logger
	reviseProvider ReviseProvider
	reviseManager  ReviseManager
}

func NewRevise(log *slog.Logger, reviseProvider ReviseProvider, reviseManager ReviseManager) *Revise {
	return &Revise{
		log:            log,
		reviseProvider: reviseProvider,
		reviseManager:  reviseManager,
	}
}

func (r *Revise) Get(ctx context.Context, id string) (domain.ReviseItem, error) {
	const op = "service.revise.get"
	panic("not implemented") // TODO: Implement
}

func (r *Revise) List(ctx context.Context, userID string) ([]domain.ReviseItem, domain.PaginationMetadata, error) {
	const op = "service.revise.list"
	panic("not implemented") // TODO: Implement
}

func (r *Revise) Create(ctx context.Context, revise domain.ReviseItem) (string, error) {
	const op = "service.revise.create"
	panic("not implemented") // TODO: Implement
}

func (r *Revise) Update(ctx context.Context, revise domain.ReviseItem) error {
	const op = "service.revise.update"
	panic("not implemented") // TODO: Implement
}

func (r *Revise) Delete(ctx context.Context, id string) error {
	const op = "service.revise.delete"
	panic("not implemented") // TODO: Implement
}
