package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aurora/aurora-backend/internal/domain"
)

type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{pool: pool}
}

func (r *SubscriptionRepo) List(ctx context.Context) ([]domain.Subscription, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.id, s.user_id, u.username, s.token, s.url, s.format,
			s.enabled, s.last_request_at, COALESCE(s.last_user_agent,''), s.created_at
		FROM subscriptions s
		JOIN users u ON u.id = s.user_id
		ORDER BY s.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		s := domain.Subscription{}
		var lastReq *time.Time
		err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Token, &s.URL,
			&s.Format, &s.Enabled, &lastReq, &s.LastUserAgent, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		s.LastRequestAt = lastReq
		subs = append(subs, s)
	}
	return subs, nil
}

func (r *SubscriptionRepo) FindByID(ctx context.Context, id string) (*domain.Subscription, error) {
	s := &domain.Subscription{}
	var lastReq *time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT s.id, s.user_id, u.username, s.token, s.url, s.format,
			s.enabled, s.last_request_at, COALESCE(s.last_user_agent,''), s.created_at
		FROM subscriptions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = $1`, id,
	).Scan(&s.ID, &s.UserID, &s.Username, &s.Token, &s.URL,
		&s.Format, &s.Enabled, &lastReq, &s.LastUserAgent, &s.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.LastRequestAt = lastReq
	return s, nil
}

func (r *SubscriptionRepo) FindByToken(ctx context.Context, token string) (*domain.Subscription, error) {
	s := &domain.Subscription{}
	var lastReq *time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT s.id, s.user_id, u.username, s.token, s.url, s.format,
			s.enabled, s.last_request_at, COALESCE(s.last_user_agent,''), s.created_at
		FROM subscriptions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token = $1`, token,
	).Scan(&s.ID, &s.UserID, &s.Username, &s.Token, &s.URL,
		&s.Format, &s.Enabled, &lastReq, &s.LastUserAgent, &s.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.LastRequestAt = lastReq
	return s, nil
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO subscriptions (user_id, token, url, format, enabled)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at`,
		sub.UserID, sub.Token, sub.URL, sub.Format, sub.Enabled,
	).Scan(&sub.ID, &sub.CreatedAt)
}

func (r *SubscriptionRepo) Toggle(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE subscriptions SET enabled = NOT enabled WHERE id = $1`, id)
	return err
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	return err
}

func (r *SubscriptionRepo) UpdateLastRequest(ctx context.Context, id, userAgent string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE subscriptions SET last_request_at = now(), last_user_agent = $2 WHERE id = $1`,
		id, userAgent)
	return err
}
