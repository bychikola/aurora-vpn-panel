package service

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/repository"
)

type DashboardService struct {
	userRepo   repository.UserRepository
	nodeRepo   repository.NodeRepository
	trafficRepo repository.TrafficRepository
	redis      *redis.Client
}

func NewDashboardService(
	userRepo repository.UserRepository,
	nodeRepo repository.NodeRepository,
	trafficRepo repository.TrafficRepository,
	redis *redis.Client,
) *DashboardService {
	return &DashboardService{
		userRepo:    userRepo,
		nodeRepo:    nodeRepo,
		trafficRepo: trafficRepo,
		redis:       redis,
	}
}

func (s *DashboardService) GetStats(ctx context.Context) (*dto.DashboardResponse, error) {
	// Try Redis cache first
	cached, err := s.redis.Get(ctx, "dashboard:stats").Result()
	if err == nil && cached != "" {
		var stats dto.DashboardResponse
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return &stats, nil
		}
	}

	// Build from DB
	active, _, _, _ := s.userRepo.CountByStatus(ctx)
	totalUsers := active // simplified
	onlineNodes, _, _, _ := s.nodeRepo.CountByStatus(ctx)
	totalNodes := onlineNodes // simplified

	protocolDist, _ := s.userRepo.ProtocolDistribution(ctx)
	userGrowth, _ := s.userRepo.UserGrowth(ctx, 30)
	trafficHistory, _ := s.trafficRepo.TrafficHistory(ctx, 24)

	stats := &dto.DashboardResponse{
		TotalUsers:        totalUsers,
		ActiveUsers:       active,
		TotalNodes:        totalNodes,
		OnlineNodes:       onlineNodes,
		TotalTrafficUp:    0,
		TotalTrafficDown:  0,
		ActiveConnections: 0,
	}

	// Convert protocol distribution
	stats.ProtocolDistribution = make([]dto.ProtocolCount, len(protocolDist))
	for i, pc := range protocolDist {
		stats.ProtocolDistribution[i] = dto.ProtocolCount{
			Protocol: pc.Protocol,
			Count:    pc.Count,
		}
	}

	// Convert traffic history
	stats.TrafficHistory = make([]dto.TrafficPoint, len(trafficHistory))
	for i, tp := range trafficHistory {
		stats.TrafficHistory[i] = dto.TrafficPoint{
			Timestamp: tp.Timestamp,
			Upload:    tp.Upload,
			Download:  tp.Download,
		}
	}

	// Convert user growth
	stats.UserGrowth = make([]dto.GrowthPoint, len(userGrowth))
	for i, gp := range userGrowth {
		stats.UserGrowth[i] = dto.GrowthPoint{
			Date:  gp.Date,
			Count: gp.Count,
		}
	}

	// Cache for 15 seconds
	data, _ := json.Marshal(stats)
	s.redis.Set(ctx, "dashboard:stats", data, 15*1e9)

	return stats, nil
}
