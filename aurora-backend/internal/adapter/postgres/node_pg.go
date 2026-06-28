package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aurora/aurora-backend/internal/domain"
)

type NodeRepo struct {
	pool *pgxpool.Pool
}

func NewNodeRepo(pool *pgxpool.Pool) *NodeRepo {
	return &NodeRepo{pool: pool}
}

func (r *NodeRepo) List(ctx context.Context) ([]domain.Node, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT n.id, n.name, n.host, n.port, n.api_port, n.api_key,
		        n.status, n.version, n.location, n.weight, n.enabled,
		        n.last_ping_at, n.metrics, n.created_at, n.updated_at,
		        COALESCE((SELECT COUNT(*) FROM inbounds WHERE node_id = n.id), 0) as inbound_count,
		        COALESCE((SELECT COUNT(*) FROM user_inbounds ui
		         JOIN inbounds i ON i.id = ui.inbound_id
		         WHERE i.node_id = n.id), 0) as user_count
		 FROM nodes n
		 ORDER BY n.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		n := domain.Node{}
		var lastPingAt *time.Time
		var userCount, inboundCount int
		err := rows.Scan(
			&n.ID, &n.Name, &n.Host, &n.Port, &n.APIPort, &n.APIKey,
			&n.Status, &n.Version, &n.Location, &n.Weight, &n.Enabled,
			&lastPingAt, &n.Metrics, &n.CreatedAt, &n.UpdatedAt,
			&inboundCount, &userCount,
		)
		if err != nil {
			return nil, err
		}
		n.LastPingAt = lastPingAt
		// Attach counts via a separate field isn't in domain.Node, we'll handle in service layer
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (r *NodeRepo) FindByID(ctx context.Context, id string) (*domain.Node, error) {
	n := &domain.Node{}
	var lastPingAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, name, host, port, api_port, api_key,
		        status, version, location, weight, enabled,
		        last_ping_at, metrics, created_at, updated_at
		 FROM nodes WHERE id = $1`, id,
	).Scan(&n.ID, &n.Name, &n.Host, &n.Port, &n.APIPort, &n.APIKey,
		&n.Status, &n.Version, &n.Location, &n.Weight, &n.Enabled,
		&lastPingAt, &n.Metrics, &n.CreatedAt, &n.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	n.LastPingAt = lastPingAt
	return n, nil
}

func (r *NodeRepo) Create(ctx context.Context, node *domain.Node) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO nodes (name, host, port, api_port, api_key, location, weight, enabled)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, status, created_at, updated_at`,
		node.Name, node.Host, node.Port, node.APIPort, node.APIKey,
		node.Location, node.Weight, node.Enabled,
	).Scan(&node.ID, &node.Status, &node.CreatedAt, &node.UpdatedAt)
}

func (r *NodeRepo) Update(ctx context.Context, node *domain.Node) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE nodes SET name=$2, host=$3, port=$4, api_port=$5,
		        api_key=$6, location=$7, weight=$8, enabled=$9,
		        updated_at=now()
		 WHERE id=$1`,
		node.ID, node.Name, node.Host, node.Port, node.APIPort,
		node.APIKey, node.Location, node.Weight, node.Enabled)
	return err
}

func (r *NodeRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM nodes WHERE id = $1`, id)
	return err
}

func (r *NodeRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE nodes SET status = $2, updated_at = now() WHERE id = $1`, id, status)
	return err
}

func (r *NodeRepo) UpdateMetrics(ctx context.Context, id, metrics string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE nodes SET metrics = $2, updated_at = now() WHERE id = $1`, id, metrics)
	return err
}

func (r *NodeRepo) UpdateLastPing(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE nodes SET last_ping_at = now() WHERE id = $1`, id)
	return err
}

func (r *NodeRepo) CountByStatus(ctx context.Context) (int, int, int, error) {
	var online, offline, degraded int
	err := r.pool.QueryRow(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE status = 'online'),
			COUNT(*) FILTER (WHERE status = 'offline'),
			COUNT(*) FILTER (WHERE status = 'degraded')
		 FROM nodes`).Scan(&online, &offline, &degraded)
	return online, offline, degraded, err
}
