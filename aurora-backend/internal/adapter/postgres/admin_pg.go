package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aurora/aurora-backend/internal/domain"
)

type AdminRepo struct {
	pool *pgxpool.Pool
}

func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{pool: pool}
}

func (r *AdminRepo) FindByUsername(ctx context.Context, username string) (*domain.Admin, error) {
	admin := &domain.Admin{}
	var lockedUntil *time.Time
	var lastLoginAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, role, failed_attempts,
		        locked_until, last_login_at, created_at, updated_at
		 FROM admins WHERE username = $1`, username,
	).Scan(&admin.ID, &admin.Username, &admin.PasswordHash, &admin.Role,
		&admin.FailedAttempts, &lockedUntil, &lastLoginAt,
		&admin.CreatedAt, &admin.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	admin.LockedUntil = lockedUntil
	admin.LastLoginAt = lastLoginAt
	return admin, nil
}

func (r *AdminRepo) FindByID(ctx context.Context, id string) (*domain.Admin, error) {
	admin := &domain.Admin{}
	var lockedUntil *time.Time
	var lastLoginAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, role, failed_attempts,
		        locked_until, last_login_at, created_at, updated_at
		 FROM admins WHERE id = $1`, id,
	).Scan(&admin.ID, &admin.Username, &admin.PasswordHash, &admin.Role,
		&admin.FailedAttempts, &lockedUntil, &lastLoginAt,
		&admin.CreatedAt, &admin.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	admin.LockedUntil = lockedUntil
	admin.LastLoginAt = lastLoginAt
	return admin, nil
}

func (r *AdminRepo) UpdateLoginAttempts(ctx context.Context, id string, attempts int, lockedUntil *string) error {
	var lu any
	if lockedUntil != nil {
		t, err := time.Parse(time.RFC3339, *lockedUntil)
		if err != nil {
			return err
		}
		lu = t
	}

	_, err := r.pool.Exec(ctx,
		`UPDATE admins SET failed_attempts = $2, locked_until = $3, updated_at = now()
		 WHERE id = $1`, id, attempts, lu)
	return err
}

func (r *AdminRepo) UpdateLastLogin(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE admins SET last_login_at = now(), failed_attempts = 0,
		        locked_until = NULL, updated_at = now()
		 WHERE id = $1`, id)
	return err
}

func (r *AdminRepo) Create(ctx context.Context, admin *domain.Admin) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO admins (username, password_hash, role)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		admin.Username, admin.PasswordHash, admin.Role,
	).Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt)
}
