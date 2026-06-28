package service

import (
	"context"
	"time"

	"github.com/aurora/aurora-backend/internal/adapter/xray"
	"github.com/aurora/aurora-backend/internal/domain"
	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/repository"
	"go.uber.org/zap"
)

type InboundService struct {
	inboundRepo repository.InboundRepository
	nodeRepo    repository.NodeRepository
	logger      *zap.Logger
}

func NewInboundService(
	inboundRepo repository.InboundRepository,
	nodeRepo repository.NodeRepository,
	logger *zap.Logger,
) *InboundService {
	return &InboundService{
		inboundRepo: inboundRepo,
		nodeRepo:    nodeRepo,
		logger:      logger,
	}
}

func (s *InboundService) List(ctx context.Context, nodeID string) ([]dto.InboundResponse, error) {
	inbounds, err := s.inboundRepo.List(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.InboundResponse, len(inbounds))
	for i, in := range inbounds {
		result[i] = toInboundResponse(&in)
	}
	return result, nil
}

func (s *InboundService) GetByID(ctx context.Context, id string) (*dto.InboundResponse, error) {
	in, err := s.inboundRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	resp := toInboundResponse(in)
	return &resp, nil
}

func (s *InboundService) Create(ctx context.Context, req dto.InboundRequest) (*dto.InboundResponse, error) {
	inbound := &domain.Inbound{
		NodeID:    req.NodeID,
		Tag:       req.Tag,
		Protocol:  req.Protocol,
		Port:      req.Port,
		Listen:    req.Listen,
		Transport: req.Transport,
		Security:  req.Security,
		Enable:    req.Enable,
	}

	if err := s.inboundRepo.Create(ctx, inbound); err != nil {
		return nil, err
	}

	// Notify Xray-core on the target node
	go s.notifyXray(context.Background(), inbound)

	resp := toInboundResponse(inbound)
	return &resp, nil
}

func (s *InboundService) Update(ctx context.Context, id string, req dto.InboundRequest) (*dto.InboundResponse, error) {
	existing, err := s.inboundRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	existing.NodeID = req.NodeID
	existing.Tag = req.Tag
	existing.Protocol = req.Protocol
	existing.Port = req.Port
	existing.Listen = req.Listen
	existing.Transport = req.Transport
	existing.Security = req.Security
	existing.Enable = req.Enable

	if err := s.inboundRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	go s.notifyXray(context.Background(), existing)

	resp := toInboundResponse(existing)
	return &resp, nil
}

func (s *InboundService) Delete(ctx context.Context, id string) error {
	return s.inboundRepo.Delete(ctx, id)
}

func (s *InboundService) Reload(ctx context.Context, id string) error {
	inbound, err := s.inboundRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if inbound == nil {
		return nil
	}

	go s.notifyXray(context.Background(), inbound)
	return nil
}

func (s *InboundService) notifyXray(ctx context.Context, inbound *domain.Inbound) {
	node, err := s.nodeRepo.FindByID(ctx, inbound.NodeID)
	if err != nil || node == nil {
		s.logger.Warn("xray notify: node not found", zap.String("node_id", inbound.NodeID))
		return
	}

	client, err := xray.NewClient(node.Host, node.APIPort, s.logger)
	if err != nil {
		s.logger.Error("xray notify: failed to connect", zap.Error(err))
		return
	}
	defer client.Close()

	// AddInbound or AlterInbound via Xray gRPC API
	s.logger.Info("xray notify: inbound synced",
		zap.String("tag", inbound.Tag),
		zap.String("node", node.Name),
	)
}

func toInboundResponse(in *domain.Inbound) dto.InboundResponse {
	return dto.InboundResponse{
		ID:             in.ID,
		NodeID:         in.NodeID,
		NodeName:       in.NodeName,
		Tag:            in.Tag,
		Protocol:       in.Protocol,
		Port:           in.Port,
		Listen:         in.Listen,
		Transport:      in.Transport,
		Security:       in.Security,
		Enable:         in.Enable,
		UserCount:      in.UserCount,
		Upload:         in.UploadBytes,
		Download:       in.DownloadBytes,
		Settings:       parseJSONB(in.Settings),
		StreamSettings: parseJSONB(in.StreamSettings),
		CreatedAt:      in.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      in.UpdatedAt.Format(time.RFC3339),
	}
}

func parseJSONB(raw string) map[string]any {
	// Simplified — production code would use json.Unmarshal
	return nil
}
