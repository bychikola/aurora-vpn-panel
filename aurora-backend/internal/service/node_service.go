package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aurora/aurora-backend/internal/domain"
	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/repository"

	auroraCrypto "github.com/aurora/aurora-backend/internal/pkg/crypto"
)

type NodeService struct {
	nodeRepo    repository.NodeRepository
	inboundRepo repository.InboundRepository
	aes         *auroraCrypto.AESGCM
}

func NewNodeService(
	nodeRepo repository.NodeRepository,
	inboundRepo repository.InboundRepository,
	aes *auroraCrypto.AESGCM,
) *NodeService {
	return &NodeService{
		nodeRepo:    nodeRepo,
		inboundRepo: inboundRepo,
		aes:         aes,
	}
}

func (s *NodeService) List(ctx context.Context) ([]dto.NodeResponse, error) {
	nodes, err := s.nodeRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.NodeResponse, len(nodes))
	for i, n := range nodes {
		result[i] = toNodeResponse(&n)
	}
	return result, nil
}

func (s *NodeService) GetByID(ctx context.Context, id string) (*dto.NodeResponse, error) {
	n, err := s.nodeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	resp := toNodeResponse(n)
	return &resp, nil
}

func (s *NodeService) Create(ctx context.Context, req dto.NodeRequest) (*dto.NodeResponse, error) {
	// Encrypt API key
	encryptedKey, err := s.aes.Encrypt(req.APIKey)
	if err != nil {
		return nil, err
	}

	node := &domain.Node{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		APIPort:  req.APIPort,
		APIKey:   encryptedKey,
		Location: req.Location,
		Weight:   req.Weight,
		Enabled:  true,
	}

	if err := s.nodeRepo.Create(ctx, node); err != nil {
		return nil, err
	}

	resp := toNodeResponse(node)
	return &resp, nil
}

func (s *NodeService) Update(ctx context.Context, id string, req dto.NodeRequest) (*dto.NodeResponse, error) {
	existing, err := s.nodeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	encryptedKey, err := s.aes.Encrypt(req.APIKey)
	if err != nil {
		return nil, err
	}

	existing.Name = req.Name
	existing.Host = req.Host
	existing.Port = req.Port
	existing.APIPort = req.APIPort
	existing.APIKey = encryptedKey
	existing.Location = req.Location
	existing.Weight = req.Weight

	if err := s.nodeRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	resp := toNodeResponse(existing)
	return &resp, nil
}

func (s *NodeService) Delete(ctx context.Context, id string) error {
	return s.nodeRepo.Delete(ctx, id)
}

func toNodeResponse(n *domain.Node) dto.NodeResponse {
	var metrics map[string]float64
	if n.Metrics != "" {
		_ = json.Unmarshal([]byte(n.Metrics), &metrics)
	}

	resp := dto.NodeResponse{
		ID:       n.ID,
		Name:     n.Name,
		Host:     n.Host,
		Port:     n.Port,
		APIPort:  n.APIPort,
		Status:   n.Status,
		Version:  n.Version,
		Location: n.Location,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}

	if metrics != nil {
		resp.CPUPercent = metrics["cpu"]
		resp.MemoryPercent = metrics["memory"]
		resp.DiskPercent = metrics["disk"]
		resp.UplinkSpeed = metrics["uplink_speed"]
		resp.DownlinkSpeed = metrics["downlink_speed"]
		resp.UplinkTotal = int64(metrics["uplink_total"])
		resp.DownlinkTotal = int64(metrics["downlink_total"])
	}

	if n.LastPingAt != nil {
		resp.LastPing = n.LastPingAt.Format(time.RFC3339)
	}

	return resp
}
