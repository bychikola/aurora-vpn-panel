package repository

import (
	"context"

	"github.com/aurora/aurora-backend/internal/domain"
)

type SettingsRepository interface {
	GetAll(ctx context.Context) ([]domain.Setting, error)
	Get(ctx context.Context, key string) (*domain.Setting, error)
	Set(ctx context.Context, key, value, description string) error
	BulkSet(ctx context.Context, settings map[string]string) error
}
