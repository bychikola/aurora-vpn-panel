package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aurora/aurora-backend/internal/domain"
)

type InboundRepo struct {
	pool *pgxpool.Pool
}

func NewInboundRepo(pool *pgxpool.Pool) *InboundRepo {
	return &InboundRepo{pool: pool}
}

func (r *InboundRepo) List(ctx context.Context, nodeID string) ([]domain.Inbound, error) {
	query := `SELECT i.id, i.node_id, i.tag, i.protocol, i.port, i.listen,
		i.transport, i.security, i.enable,
		i.settings::text, i.stream_settings::text, i.sniffing::text,
		i.user_count, i.upload_bytes, i.download_bytes,
		i.created_at, i.updated_at,
		COALESCE(n.name, '') as node_name
		FROM inbounds i
		JOIN nodes n ON n.id = i.node_id`
	args := []any{}

	if nodeID != "" {
		query += " WHERE i.node_id = $1"
		args = append(args, nodeID)
	}
	query += " ORDER BY i.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inbounds []domain.Inbound
	for rows.Next() {
		in := domain.Inbound{}
		err := rows.Scan(
			&in.ID, &in.NodeID, &in.Tag, &in.Protocol, &in.Port, &in.Listen,
			&in.Transport, &in.Security, &in.Enable,
			&in.Settings, &in.StreamSettings, &in.Sniffing,
			&in.UserCount, &in.UploadBytes, &in.DownloadBytes,
			&in.CreatedAt, &in.UpdatedAt,
			&in.NodeName,
		)
		if err != nil {
			return nil, err
		}
		inbounds = append(inbounds, in)
	}
	return inbounds, nil
}

func (r *InboundRepo) FindByID(ctx context.Context, id string) (*domain.Inbound, error) {
	in := &domain.Inbound{}
	err := r.pool.QueryRow(ctx,
		`SELECT i.id, i.node_id, i.tag, i.protocol, i.port, i.listen,
			i.transport, i.security, i.enable,
			i.settings::text, i.stream_settings::text, i.sniffing::text,
			i.user_count, i.upload_bytes, i.download_bytes,
			i.created_at, i.updated_at,
			COALESCE(n.name, '') as node_name
		FROM inbounds i
		JOIN nodes n ON n.id = i.node_id
		WHERE i.id = $1`, id,
	).Scan(&in.ID, &in.NodeID, &in.Tag, &in.Protocol, &in.Port, &in.Listen,
		&in.Transport, &in.Security, &in.Enable,
		&in.Settings, &in.StreamSettings, &in.Sniffing,
		&in.UserCount, &in.UploadBytes, &in.DownloadBytes,
		&in.CreatedAt, &in.UpdatedAt, &in.NodeName)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (r *InboundRepo) FindByTag(ctx context.Context, nodeID, tag string) (*domain.Inbound, error) {
	in := &domain.Inbound{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, node_id, tag, protocol, port, listen, transport, security, enable,
			settings::text, stream_settings::text, sniffing::text,
			user_count, upload_bytes, download_bytes, created_at, updated_at
		FROM inbounds WHERE node_id = $1 AND tag = $2`, nodeID, tag,
	).Scan(&in.ID, &in.NodeID, &in.Tag, &in.Protocol, &in.Port, &in.Listen,
		&in.Transport, &in.Security, &in.Enable,
		&in.Settings, &in.StreamSettings, &in.Sniffing,
		&in.UserCount, &in.UploadBytes, &in.DownloadBytes,
		&in.CreatedAt, &in.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (r *InboundRepo) Create(ctx context.Context, inbound *domain.Inbound) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO inbounds (node_id, tag, protocol, port, listen, transport, security, enable, settings, stream_settings, sniffing)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10::jsonb,$11::jsonb)
		RETURNING id, created_at, updated_at`,
		inbound.NodeID, inbound.Tag, inbound.Protocol, inbound.Port,
		inbound.Listen, inbound.Transport, inbound.Security, inbound.Enable,
		inbound.Settings, inbound.StreamSettings, inbound.Sniffing,
	).Scan(&inbound.ID, &inbound.CreatedAt, &inbound.UpdatedAt)
}

func (r *InboundRepo) Update(ctx context.Context, inbound *domain.Inbound) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inbounds SET node_id=$2, tag=$3, protocol=$4, port=$5, listen=$6,
		transport=$7, security=$8, enable=$9, settings=$10::jsonb,
		stream_settings=$11::jsonb, sniffing=$12::jsonb, updated_at=now()
		WHERE id=$1`,
		inbound.ID, inbound.NodeID, inbound.Tag, inbound.Protocol, inbound.Port,
		inbound.Listen, inbound.Transport, inbound.Security, inbound.Enable,
		inbound.Settings, inbound.StreamSettings, inbound.Sniffing)
	return err
}

func (r *InboundRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inbounds WHERE id = $1`, id)
	return err
}

func (r *InboundRepo) UpdateTraffic(ctx context.Context, id string, upload, download int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inbounds SET upload_bytes = upload_bytes + $2,
		download_bytes = download_bytes + $3, updated_at = now()
		WHERE id = $1`, id, upload, download)
	return err
}

func (r *InboundRepo) UpdateUserCount(ctx context.Context, id string, count int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inbounds SET user_count = $2, updated_at = now() WHERE id = $1`, id, count)
	return err
}
