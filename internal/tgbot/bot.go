package tgbot

import (
	"context"

	"github.com/ARUMANDESU/go-revise/internal/domain"
)

type Bot struct {
}

type ReviseService interface {
	Get(ctx context.Context, id string) (domain.ReviseItem, error)
	List(ctx context.Context, userID string) ([]domain.ReviseItem, domain.PaginationMetadata, error)
	Create(ctx context.Context, dto domain.CreateReviseItemDTO) (domain.ReviseItem, error)
	Update(ctx context.Context, revise domain.ReviseItem) (domain.ReviseItem, error)
	Delete(ctx context.Context, id string, userID string) (domain.ReviseItem, error)
}
