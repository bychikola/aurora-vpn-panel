package worker

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/aurora/aurora-backend/internal/adapter/xray"
	"github.com/aurora/aurora-backend/internal/repository"
)

// TrafficCollector периодически собирает статистику трафика из Xray-core нод
// и сохраняет её в Redis (мгновенный доступ) и PostgreSQL (агрегация).
type TrafficCollector struct {
	nodeRepo    repository.NodeRepository
	userRepo    repository.UserRepository
	trafficRepo repository.TrafficRepository
	redis       *redis.Client
	logger      *zap.Logger
	interval    time.Duration
}

func NewTrafficCollector(
	nodeRepo repository.NodeRepository,
	userRepo repository.UserRepository,
	trafficRepo repository.TrafficRepository,
	redis *redis.Client,
	logger *zap.Logger,
	interval time.Duration,
) *TrafficCollector {
	return &TrafficCollector{
		nodeRepo:    nodeRepo,
		userRepo:    userRepo,
		trafficRepo: trafficRepo,
		redis:       redis,
		logger:      logger,
		interval:    interval,
	}
}

// Run запускает бесконечный цикл сбора статистики.
// Должен вызываться в отдельной горутине.
func (w *TrafficCollector) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("traffic collector started", zap.Duration("interval", w.interval))

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("traffic collector stopped")
			return
		case <-ticker.C:
			w.collect(ctx)
		}
	}
}

func (w *TrafficCollector) collect(ctx context.Context) {
	// Acquire distributed lock to prevent multiple collectors
	locked, _ := w.redis.SetNX(ctx, "lock:traffic:collect", "1", w.interval).Result()
	if !locked {
		return // Another instance is collecting
	}
	defer w.redis.Del(ctx, "lock:traffic:collect")

	nodes, err := w.nodeRepo.List(ctx)
	if err != nil {
		w.logger.Error("traffic collector: failed to list nodes", zap.Error(err))
		return
	}

	for _, node := range nodes {
		if node.Status != "online" {
			continue
		}

		go func(nodeHost string, apiPort int) {
			client, err := xray.NewClient(nodeHost, apiPort, w.logger)
			if err != nil {
				w.logger.Warn("traffic collector: failed to connect to node",
					zap.String("host", nodeHost), zap.Error(err))
				return
			}
			defer client.Close()

			stats, err := client.QueryAllStats(ctx)
			if err != nil {
				w.logger.Warn("traffic collector: failed to query stats",
					zap.String("host", nodeHost), zap.Error(err))
				return
			}

			// Write to Redis (instant access for dashboard)
			for email, stat := range stats {
				if stat.Upload > 0 || stat.Download > 0 {
					w.redis.Set(ctx,
						"traffic:user:"+email+":upload", stat.Upload, 3600*time.Second)
					w.redis.Set(ctx,
						"traffic:user:"+email+":download", stat.Download, 3600*time.Second)

					// Update user traffic in PostgreSQL (batched would be better)
					_ = w.userRepo.UpdateTraffic(ctx, email, stat.Upload, stat.Download)
				}
			}

			w.logger.Debug("traffic collector: stats collected",
				zap.String("host", nodeHost),
				zap.Int("users", len(stats)),
			)
		}(node.Host, node.APIPort)
	}
}
