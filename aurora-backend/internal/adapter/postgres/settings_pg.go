package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aurora/aurora-backend/internal/domain"
)

type SettingsRepo struct {
	pool *pgxpool.Pool
}

func NewSettingsRepo(pool *pgxpool.Pool) *SettingsRepo {
	return &SettingsRepo{pool: pool}
}

func (r *SettingsRepo) GetAll(ctx context.Context) ([]domain.Setting, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT key, value, COALESCE(description,''), updated_at FROM settings ORDER BY key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []domain.Setting
	for rows.Next() {
		s := domain.Setting{}
		if err := rows.Scan(&s.Key, &s.Value, &s.Description, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SettingsRepo) Get(ctx context.Context, key string) (*domain.Setting, error) {
	s := &domain.Setting{}
	err := r.pool.QueryRow(ctx,
		`SELECT key, value, COALESCE(description,''), updated_at FROM settings WHERE key = $1`, key,
	).Scan(&s.Key, &s.Value, &s.Description, &s.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SettingsRepo) Set(ctx context.Context, key, value, description string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO settings (key, value, description, updated_at)
		VALUES ($1,$2,$3,now())
		ON CONFLICT (key) DO UPDATE SET value = $2, description = $3, updated_at = now()`,
		key, value, description)
	return err
}

func (r *SettingsRepo) BulkSet(ctx context.Context, settings map[string]string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for key, value := range settings {
		_, err := tx.Exec(ctx,
			`INSERT INTO settings (key, value, updated_at)
			VALUES ($1,$2,now())
			ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = now()`,
			key, value)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
