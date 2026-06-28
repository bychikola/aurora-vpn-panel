package repository

import (
	"context"

	"github.com/aurora/aurora-backend/internal/domain"
)

type NodeRepository interface {
	List(ctx context.Context) ([]domain.Node, error)
	FindByID(ctx context.Context, id string) (*domain.Node, error)
	Create(ctx context.Context, node *domain.Node) error
	Update(ctx context.Context, node *domain.Node) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateMetrics(ctx context.Context, id string, metrics string) error
	UpdateLastPing(ctx context.Context, id string) error
	CountByStatus(ctx context.Context) (online int, offline int, degraded int, err error)
}
