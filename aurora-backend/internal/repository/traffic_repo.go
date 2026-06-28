package repository

import (
	"context"

	"github.com/aurora/aurora-backend/internal/domain"
)

type TrafficRepository interface {
	InsertConnectionLog(ctx context.Context, log *domain.ConnectionLog) error
	InsertInboundTraffic(ctx context.Context, traffic *domain.InboundTraffic) error
	TrafficHistory(ctx context.Context, hours int) ([]domain.TrafficPoint, error)
	CleanupOldLogs(ctx context.Context, retentionDays int) error
}
