package repository

import (
	"context"

	"github.com/aurora/aurora-backend/internal/domain"
)

type AdminRepository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Admin, error)
	FindByID(ctx context.Context, id string) (*domain.Admin, error)
	UpdateLoginAttempts(ctx context.Context, id string, attempts int, lockedUntil *string) error
	UpdateLastLogin(ctx context.Context, id string) error
	Create(ctx context.Context, admin *domain.Admin) error
}
