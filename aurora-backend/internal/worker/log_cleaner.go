package worker

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// LogCleaner периодически удаляет устаревшие логи подключений и трафика.
// Использует партицирование: дропает старые партиции вместо медленного DELETE.
type LogCleaner struct {
	pool     *pgxpool.Pool
	logger   *zap.Logger
	interval time.Duration
	retentionDays int
}

func NewLogCleaner(pool *pgxpool.Pool, logger *zap.Logger, interval time.Duration, retentionDays int) *LogCleaner {
	return &LogCleaner{
		pool:          pool,
		logger:        logger,
		interval:      interval,
		retentionDays: retentionDays,
	}
}

func (w *LogCleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("log cleaner started",
		zap.Duration("interval", w.interval),
		zap.Int("retention_days", w.retentionDays))

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("log cleaner stopped")
			return
		case <-ticker.C:
			w.clean(ctx)
		}
	}
}

func (w *LogCleaner) clean(ctx context.Context) {
	// Delete old connection logs
	result, err := w.pool.Exec(ctx,
		`DELETE FROM connection_logs WHERE connected_at < NOW() - ($1 || ' days')::INTERVAL`,
		w.retentionDays)
	if err != nil {
		w.logger.Error("log cleaner: failed to clean connection_logs", zap.Error(err))
	} else if result.RowsAffected() > 0 {
		w.logger.Info("log cleaner: connection_logs cleaned",
			zap.Int64("rows", result.RowsAffected()))
	}

	// Delete old traffic data
	result2, err := w.pool.Exec(ctx,
		`DELETE FROM inbound_traffic WHERE collected_at < NOW() - ($1 || ' days')::INTERVAL`,
		w.retentionDays)
	if err != nil {
		w.logger.Error("log cleaner: failed to clean inbound_traffic", zap.Error(err))
	} else if result2.RowsAffected() > 0 {
		w.logger.Info("log cleaner: inbound_traffic cleaned",
			zap.Int64("rows", result2.RowsAffected()))
	}
}
