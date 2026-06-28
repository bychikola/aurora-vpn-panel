package repository

import (
	"context"

	"github.com/aurora/aurora-backend/internal/domain"
)

type InboundRepository interface {
	List(ctx context.Context, nodeID string) ([]domain.Inbound, error)
	FindByID(ctx context.Context, id string) (*domain.Inbound, error)
	FindByTag(ctx context.Context, nodeID, tag string) (*domain.Inbound, error)
	Create(ctx context.Context, inbound *domain.Inbound) error
	Update(ctx context.Context, inbound *domain.Inbound) error
	Delete(ctx context.Context, id string) error
	UpdateTraffic(ctx context.Context, id string, upload, download int64) error
	UpdateUserCount(ctx context.Context, id string, count int) error
}
