package repository

import (
	"context"

	"github.com/aurora/aurora-backend/internal/domain"
)

type UserRepository interface {
	List(ctx context.Context, filters UserFilter) ([]domain.User, int, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindBySubscriptionToken(ctx context.Context, token string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User, protocols []string, inboundIDs []string) error
	Update(ctx context.Context, user *domain.User, protocols []string, inboundIDs []string) error
	Delete(ctx context.Context, id string) error
	ResetTraffic(ctx context.Context, id string) error
	UpdateTraffic(ctx context.Context, email string, upload, download int64) error
	UpdateSubscriptionToken(ctx context.Context, id, token string) error
	CountByStatus(ctx context.Context) (active int, disabled int, expired int, err error)
	ProtocolDistribution(ctx context.Context) ([]domain.ProtocolCount, error)
	UserGrowth(ctx context.Context, days int) ([]domain.GrowthPoint, error)
}

type UserFilter struct {
	Search   string
	Status   string
	Protocol string
	Page     int
	PageSize int
}
