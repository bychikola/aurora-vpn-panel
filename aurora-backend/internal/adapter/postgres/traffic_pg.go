package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aurora/aurora-backend/internal/domain"
)

type TrafficRepo struct {
	pool *pgxpool.Pool
}

func NewTrafficRepo(pool *pgxpool.Pool) *TrafficRepo {
	return &TrafficRepo{pool: pool}
}

func (r *TrafficRepo) InsertConnectionLog(ctx context.Context, log *domain.ConnectionLog) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO connection_logs (user_id, email, inbound_id, node_id, ip_address, user_agent, upload_bytes, download_bytes)
		VALUES ($1,$2,$3,$4,$5::inet,$6,$7,$8)
		RETURNING id`,
		log.UserID, log.Email, log.InboundID, log.NodeID,
		log.IPAddress, log.UserAgent, log.UploadBytes, log.DownloadBytes,
	).Scan(&log.ID)
}

func (r *TrafficRepo) InsertInboundTraffic(ctx context.Context, traffic *domain.InboundTraffic) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inbound_traffic (inbound_id, user_email, upload, download)
		VALUES ($1,$2,$3,$4)`,
		traffic.InboundID, traffic.UserEmail, traffic.Upload, traffic.Download)
	return err
}

func (r *TrafficRepo) TrafficHistory(ctx context.Context, hours int) ([]domain.TrafficPoint, error) {
	// Aggregate traffic per hour from inbound_traffic
	rows, err := r.pool.Query(ctx,
		fmt.Sprintf(`SELECT
			date_trunc('hour', collected_at) as ts,
			COALESCE(SUM(upload), 0) as total_upload,
			COALESCE(SUM(download), 0) as total_download
		FROM inbound_traffic
		WHERE collected_at >= NOW() - ($1 || ' hours')::INTERVAL
		GROUP BY ts
		ORDER BY ts`, hours),
		fmt.Sprintf("%d", hours))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []domain.TrafficPoint
	for rows.Next() {
		p := domain.TrafficPoint{}
		var upload, download float64
		if err := rows.Scan(&p.Timestamp, &upload, &download); err != nil {
			return nil, err
		}
		p.Upload = upload
		p.Download = download
		points = append(points, p)
	}
	return points, nil
}

func (r *TrafficRepo) CleanupOldLogs(ctx context.Context, retentionDays int) error {
	// Drop old partitions instead of DELETE (much faster)
	// In production, you'd dynamically compute the partition name
	_, err := r.pool.Exec(ctx,
		fmt.Sprintf(`DELETE FROM connection_logs WHERE connected_at < NOW() - ($1 || ' days')::INTERVAL`,
			retentionDays),
		retentionDays)
	return err
}
