package repository

import (
	"context"

	"github.com/aurora/aurora-backend/internal/domain"
)

type SubscriptionRepository interface {
	List(ctx context.Context) ([]domain.Subscription, error)
	FindByID(ctx context.Context, id string) (*domain.Subscription, error)
	FindByToken(ctx context.Context, token string) (*domain.Subscription, error)
	Create(ctx context.Context, sub *domain.Subscription) error
	Toggle(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	UpdateLastRequest(ctx context.Context, id, userAgent string) error
}
