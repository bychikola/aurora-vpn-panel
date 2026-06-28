package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/aurora/aurora-backend/internal/adapter/xray"
)

// NodeHealthcheck периодически проверяет доступность всех нод
// и обновляет их статус в базе данных.
type NodeHealthcheck struct {
	pool     *pgxpool.Pool
	logger   *zap.Logger
	interval time.Duration
	timeout  time.Duration
}

func NewNodeHealthcheck(pool *pgxpool.Pool, logger *zap.Logger, interval, timeout time.Duration) *NodeHealthcheck {
	return &NodeHealthcheck{
		pool:     pool,
		logger:   logger,
		interval: interval,
		timeout:  timeout,
	}
}

func (w *NodeHealthcheck) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("node healthcheck started", zap.Duration("interval", w.interval))

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("node healthcheck stopped")
			return
		case <-ticker.C:
			w.check(ctx)
		}
	}
}

func (w *NodeHealthcheck) check(ctx context.Context) {
	rows, err := w.pool.Query(ctx,
		`SELECT id, host, api_port FROM nodes WHERE enabled = true`)
	if err != nil {
		w.logger.Error("healthcheck: failed to query nodes", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, host string
		var apiPort int
		if err := rows.Scan(&id, &host, &apiPort); err != nil {
			continue
		}

		go func(nodeID, nodeHost string, nodeAPIPort int) {
			checkCtx, cancel := context.WithTimeout(ctx, w.timeout)
			defer cancel()

			client, err := xray.NewClient(nodeHost, nodeAPIPort, w.logger)
			if err != nil {
				w.updateStatus(checkCtx, nodeID, "offline")
				return
			}
			defer client.Close()

			if err := client.HealthCheck(checkCtx); err != nil {
				w.updateStatus(checkCtx, nodeID, "degraded")
				w.logger.Warn("healthcheck: node unreachable",
					zap.String("host", nodeHost), zap.Error(err))
				return
			}

			w.updateStatus(checkCtx, nodeID, "online")
			w.pool.Exec(checkCtx,
				`UPDATE nodes SET last_ping_at = now() WHERE id = $1`, nodeID)
		}(id, host, apiPort)
	}
}

func (w *NodeHealthcheck) updateStatus(ctx context.Context, nodeID, status string) {
	_, err := w.pool.Exec(ctx,
		`UPDATE nodes SET status = $2, updated_at = now() WHERE id = $1 AND status != $2`,
		nodeID, status)
	if err != nil {
		w.logger.Error("healthcheck: failed to update status",
			zap.String("node_id", nodeID), zap.Error(err))
	}
}

// Metrics helper — unused for now, will be used by node agent sync
func metricsToJSON(cpu, mem, disk, upSpeed, downSpeed float64, upTotal, downTotal int64) string {
	m := map[string]any{
		"cpu":           cpu,
		"memory":        mem,
		"disk":          disk,
		"uplink_speed":  upSpeed,
		"downlink_speed": downSpeed,
		"uplink_total":  upTotal,
		"downlink_total": downTotal,
	}
	data, _ := json.Marshal(m)
	return string(data)
}
