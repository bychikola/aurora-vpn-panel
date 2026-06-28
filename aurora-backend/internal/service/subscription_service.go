package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aurora/aurora-backend/internal/domain"
	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/pkg/subscription"
	"github.com/aurora/aurora-backend/internal/repository"
)

type SubscriptionService struct {
	subRepo    repository.SubscriptionRepository
	userRepo   repository.UserRepository
	inboundRepo repository.InboundRepository
	nodeRepo   repository.NodeRepository
}

func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	userRepo repository.UserRepository,
	inboundRepo repository.InboundRepository,
	nodeRepo repository.NodeRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subRepo:    subRepo,
		userRepo:   userRepo,
		inboundRepo: inboundRepo,
		nodeRepo:   nodeRepo,
	}
}

func (s *SubscriptionService) List(ctx context.Context) ([]dto.SubscriptionResponse, error) {
	subs, err := s.subRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.SubscriptionResponse, len(subs))
	for i, sub := range subs {
		result[i] = toSubscriptionResponse(&sub)
	}
	return result, nil
}

func (s *SubscriptionService) Toggle(ctx context.Context, id string) (*dto.SubscriptionResponse, error) {
	if err := s.subRepo.Toggle(ctx, id); err != nil {
		return nil, err
	}
	// Get updated
	subs, _ := s.subRepo.List(ctx)
	for _, sub := range subs {
		if sub.ID == id {
			resp := toSubscriptionResponse(&sub)
			return &resp, nil
		}
	}
	return nil, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	return s.subRepo.Delete(ctx, id)
}

// GenerateSubscription генерирует конфиг подписки для клиента
func (s *SubscriptionService) GenerateSubscription(
	ctx context.Context, token, format, userAgent string,
) ([]byte, string, error) {
	// 1. Find user by subscription token
	user, err := s.userRepo.FindBySubscriptionToken(ctx, token)
	if err != nil {
		return nil, "", fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, "", fmt.Errorf("invalid subscription token")
	}

	// 2. Validate user status
	if user.Status != "active" {
		return nil, "", fmt.Errorf("subscription not active")
	}
	if user.ExpireAt != nil && user.ExpireAt.Before(time.Now()) {
		return nil, "", fmt.Errorf("subscription expired")
	}
	if user.TrafficLimitBytes > 0 && user.TrafficUsedBytes >= user.TrafficLimitBytes {
		return nil, "", fmt.Errorf("traffic limit exceeded")
	}

	// 3. Build proxy list from user's inbounds
	var proxies []subscription.ProxyConfig
	for _, inboundID := range user.InboundIDs {
		inbound, err := s.inboundRepo.FindByID(ctx, inboundID)
		if err != nil || inbound == nil || !inbound.Enable {
			continue
		}
		node, err := s.nodeRepo.FindByID(ctx, inbound.NodeID)
		if err != nil || node == nil || !node.Enabled {
			continue
		}

		proxy := subscription.ProxyConfig{
			Protocol:   inbound.Protocol,
			Name:       fmt.Sprintf("%s | %s", inbound.Tag, node.Name),
			Address:    node.Host,
			Port:       inbound.Port,
			Security:   inbound.Security,
			Transport:  inbound.Transport,
		}
		// In production: parse inbound.Settings for UUID/password/flow/SNI etc.
		proxies = append(proxies, proxy)
	}

	if len(proxies) == 0 {
		return nil, "", fmt.Errorf("no active proxies configured")
	}

	// 4. Detect format from User-Agent if not specified
	subFormat := subscription.FormatBase64
	if format != "" {
		subFormat = subscription.Format(format)
	} else {
		subFormat = subscription.DetectFormat(userAgent)
	}

	// 5. Generate config
	var config string
	switch subFormat {
	case subscription.FormatClash:
		config, err = subscription.GenerateClash(proxies, "AURORA")
		contentType := "text/plain; charset=utf-8"
		return []byte(config), contentType, err
	case subscription.FormatSingbox:
		config, err = subscription.GenerateSingbox(proxies)
		contentType := "application/json; charset=utf-8"
		return []byte(config), contentType, err
	default:
		config, err = subscription.GenerateBase64(proxies)
		contentType := "text/plain; charset=utf-8"
		return []byte(config), contentType, err
	}
}

func toSubscriptionResponse(sub *domain.Subscription) dto.SubscriptionResponse {
	return dto.SubscriptionResponse{
		ID:            sub.ID,
		UserID:        sub.UserID,
		Username:      sub.Username,
		Token:         sub.Token,
		URL:           sub.URL,
		Format:        sub.Format,
		Enabled:       sub.Enabled,
		LastRequestAt: dto.TimePtr(nilTime(sub.LastRequestAt)),
		LastUserAgent: sub.LastUserAgent,
		CreatedAt:     sub.CreatedAt.Format(time.RFC3339),
	}
}

func nilTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
