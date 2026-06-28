package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aurora/aurora-backend/internal/domain"
	"github.com/aurora/aurora-backend/internal/repository"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) List(ctx context.Context, f repository.UserFilter) ([]domain.User, int, error) {
	var conditions []string
	args := []any{}
	argIdx := 1

	if f.Search != "" {
		conditions = append(conditions,
			fmt.Sprintf("(u.username ILIKE $%d OR u.email ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+f.Search+"%")
		argIdx++
	}
	if f.Status != "" && f.Status != "all" {
		conditions = append(conditions, fmt.Sprintf("u.status = $%d", argIdx))
		args = append(args, f.Status)
		argIdx++
	}
	if f.Protocol != "" && f.Protocol != "all" {
		conditions = append(conditions,
			fmt.Sprintf(`EXISTS(SELECT 1 FROM user_protocols up
			 WHERE up.user_id = u.id AND up.protocol = $%d)`, argIdx))
		args = append(args, f.Protocol)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users u %s", where)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Paginated query
	offset := (f.Page - 1) * f.PageSize
	dataQuery := fmt.Sprintf(
		`SELECT u.id, u.username, u.email, u.status,
		        u.traffic_limit_bytes, u.traffic_used_bytes,
		        u.expire_at, u.max_ips, u.concurrent_ips,
		        u.subscription_token, u.notes, u.last_seen_at,
		        u.created_at, u.updated_at
		 FROM users u %s
		 ORDER BY u.created_at DESC
		 LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1,
	)
	args = append(args, f.PageSize, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		u := domain.User{}
		var expireAt, lastSeenAt *time.Time
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.Status,
			&u.TrafficLimitBytes, &u.TrafficUsedBytes,
			&expireAt, &u.MaxIPs, &u.ConcurrentIPs,
			&u.SubscriptionToken, &u.Notes, &lastSeenAt,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		u.ExpireAt = expireAt
		u.LastSeenAt = lastSeenAt
		users = append(users, u)
	}

	// Load protocols and inbound IDs for each user
	for i := range users {
		users[i].Protocols, _ = r.getProtocols(ctx, users[i].ID)
		users[i].InboundIDs, _ = r.getInboundIDs(ctx, users[i].ID)
	}

	return users, total, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	var expireAt, lastSeenAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, status, traffic_limit_bytes, traffic_used_bytes,
		        expire_at, max_ips, concurrent_ips, subscription_token, notes,
		        last_seen_at, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Status,
		&u.TrafficLimitBytes, &u.TrafficUsedBytes,
		&expireAt, &u.MaxIPs, &u.ConcurrentIPs,
		&u.SubscriptionToken, &u.Notes, &lastSeenAt,
		&u.CreatedAt, &u.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.ExpireAt = expireAt
	u.LastSeenAt = lastSeenAt

	u.Protocols, _ = r.getProtocols(ctx, u.ID)
	u.InboundIDs, _ = r.getInboundIDs(ctx, u.ID)

	return u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	var expireAt, lastSeenAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, status, traffic_limit_bytes, traffic_used_bytes,
		        expire_at, max_ips, concurrent_ips, subscription_token, notes,
		        last_seen_at, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Status,
		&u.TrafficLimitBytes, &u.TrafficUsedBytes,
		&expireAt, &u.MaxIPs, &u.ConcurrentIPs,
		&u.SubscriptionToken, &u.Notes, &lastSeenAt,
		&u.CreatedAt, &u.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.ExpireAt = expireAt
	u.LastSeenAt = lastSeenAt
	return u, nil
}

func (r *UserRepo) FindBySubscriptionToken(ctx context.Context, token string) (*domain.User, error) {
	u := &domain.User{}
	var expireAt, lastSeenAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, status, traffic_limit_bytes, traffic_used_bytes,
		        expire_at, max_ips, concurrent_ips, subscription_token, notes,
		        last_seen_at, created_at, updated_at
		 FROM users WHERE subscription_token = $1`, token,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Status,
		&u.TrafficLimitBytes, &u.TrafficUsedBytes,
		&expireAt, &u.MaxIPs, &u.ConcurrentIPs,
		&u.SubscriptionToken, &u.Notes, &lastSeenAt,
		&u.CreatedAt, &u.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.ExpireAt = expireAt
	u.LastSeenAt = lastSeenAt
	return u, nil
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User, protocols []string, inboundIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO users (username, email, status, traffic_limit_bytes,
		        expire_at, max_ips, subscription_token, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, created_at, updated_at`,
		user.Username, user.Email, user.Status, user.TrafficLimitBytes,
		user.ExpireAt, user.MaxIPs, user.SubscriptionToken, user.Notes,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	// Protocols
	for _, p := range protocols {
		_, err = tx.Exec(ctx,
			`INSERT INTO user_protocols (user_id, protocol) VALUES ($1, $2)`,
			user.ID, p)
		if err != nil {
			return err
		}
	}

	// Inbound links
	for _, inboundID := range inboundIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO user_inbounds (user_id, inbound_id) VALUES ($1, $2)`,
			user.ID, inboundID)
		if err != nil {
			return err
		}
		// Update inbound user count
		_, _ = tx.Exec(ctx,
			`UPDATE inbounds SET user_count = user_count + 1 WHERE id = $1`, inboundID)
	}

	return tx.Commit(ctx)
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User, protocols []string, inboundIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE users SET username=$2, email=$3, status=$4,
		        traffic_limit_bytes=$5, expire_at=$6, max_ips=$7,
		        notes=$8, updated_at=now()
		 WHERE id=$1`,
		user.ID, user.Username, user.Email, user.Status,
		user.TrafficLimitBytes, user.ExpireAt, user.MaxIPs,
		user.Notes)
	if err != nil {
		return err
	}

	// Replace protocols
	_, _ = tx.Exec(ctx, `DELETE FROM user_protocols WHERE user_id = $1`, user.ID)
	for _, p := range protocols {
		_, _ = tx.Exec(ctx,
			`INSERT INTO user_protocols (user_id, protocol) VALUES ($1, $2)`,
			user.ID, p)
	}

	// Replace inbound links (decrement old, increment new)
	oldInbounds, _ := r.getInboundIDs(ctx, user.ID)
	for _, oid := range oldInbounds {
		_, _ = tx.Exec(ctx,
			`UPDATE inbounds SET user_count = GREATEST(user_count - 1, 0) WHERE id = $1`, oid)
	}
	_, _ = tx.Exec(ctx, `DELETE FROM user_inbounds WHERE user_id = $1`, user.ID)
	for _, iid := range inboundIDs {
		_, _ = tx.Exec(ctx,
			`INSERT INTO user_inbounds (user_id, inbound_id) VALUES ($1, $2)`, user.ID, iid)
		_, _ = tx.Exec(ctx,
			`UPDATE inbounds SET user_count = user_count + 1 WHERE id = $1`, iid)
	}

	return tx.Commit(ctx)
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Decrement inbound user counts
	inboundIDs, _ := r.getInboundIDs(ctx, id)
	for _, iid := range inboundIDs {
		_, _ = tx.Exec(ctx,
			`UPDATE inbounds SET user_count = GREATEST(user_count - 1, 0) WHERE id = $1`, iid)
	}

	_, err = tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserRepo) ResetTraffic(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET traffic_used_bytes = 0, updated_at = now() WHERE id = $1`, id)
	return err
}

func (r *UserRepo) UpdateTraffic(ctx context.Context, email string, upload, download int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET traffic_used_bytes = traffic_used_bytes + $2 + $3,
		        updated_at = now()
		 WHERE email = $1`, email, upload, download)
	return err
}

func (r *UserRepo) UpdateSubscriptionToken(ctx context.Context, id, token string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET subscription_token = $2, updated_at = now() WHERE id = $1`,
		id, token)
	return err
}

func (r *UserRepo) CountByStatus(ctx context.Context) (int, int, int, error) {
	var active, disabled, expired int
	err := r.pool.QueryRow(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE status = 'active'),
			COUNT(*) FILTER (WHERE status = 'disabled'),
			COUNT(*) FILTER (WHERE status = 'expired')
		 FROM users`).Scan(&active, &disabled, &expired)
	return active, disabled, expired, err
}

func (r *UserRepo) ProtocolDistribution(ctx context.Context) ([]domain.ProtocolCount, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT protocol, COUNT(*) as count
		 FROM user_protocols
		 GROUP BY protocol ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ProtocolCount
	for rows.Next() {
		var pc domain.ProtocolCount
		if err := rows.Scan(&pc.Protocol, &pc.Count); err != nil {
			return nil, err
		}
		result = append(result, pc)
	}
	return result, nil
}

func (r *UserRepo) UserGrowth(ctx context.Context, days int) ([]domain.GrowthPoint, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DATE(created_at) as date, COUNT(*) as count
		 FROM users
		 WHERE created_at >= NOW() - ($1 || ' days')::INTERVAL
		 GROUP BY DATE(created_at)
		 ORDER BY date`, fmt.Sprintf("%d", days))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.GrowthPoint
	for rows.Next() {
		var gp domain.GrowthPoint
		if err := rows.Scan(&gp.Date, &gp.Count); err != nil {
			return nil, err
		}
		result = append(result, gp)
	}
	return result, nil
}

// ─── Helpers ───

func (r *UserRepo) getProtocols(ctx context.Context, userID string) ([]domain.Protocol, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, protocol FROM user_protocols WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var protocols []domain.Protocol
	for rows.Next() {
		var p domain.Protocol
		if err := rows.Scan(&p.UserID, &p.Protocol); err != nil {
			return nil, err
		}
		protocols = append(protocols, p)
	}
	return protocols, nil
}

func (r *UserRepo) getInboundIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT inbound_id FROM user_inbounds WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
