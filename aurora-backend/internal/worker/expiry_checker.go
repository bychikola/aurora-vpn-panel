package worker

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// ExpiryChecker периодически проверяет пользователей с истёкшим сроком действия
// и автоматически переводит их в статус expired.
type ExpiryChecker struct {
	pool     *pgxpool.Pool
	logger   *zap.Logger
	interval time.Duration
}

func NewExpiryChecker(pool *pgxpool.Pool, logger *zap.Logger, interval time.Duration) *ExpiryChecker {
	return &ExpiryChecker{pool: pool, logger: logger, interval: interval}
}

func (w *ExpiryChecker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("expiry checker started", zap.Duration("interval", w.interval))

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("expiry checker stopped")
			return
		case <-ticker.C:
			w.check(ctx)
		}
	}
}

func (w *ExpiryChecker) check(ctx context.Context) {
	result, err := w.pool.Exec(ctx,
		`UPDATE users SET status = 'expired', updated_at = now()
		WHERE status = 'active' AND expire_at IS NOT NULL AND expire_at < now()`)
	if err != nil {
		w.logger.Error("expiry checker: failed to update", zap.Error(err))
		return
	}

	if result.RowsAffected() > 0 {
		w.logger.Info("expiry checker: users expired",
			zap.Int64("count", result.RowsAffected()))
	}

	// Also unban IPs whose ban has expired
	_, _ = w.pool.Exec(ctx,
		`DELETE FROM banned_ips WHERE banned_until < now()`)
}
