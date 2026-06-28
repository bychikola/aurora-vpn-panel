package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/aurora/aurora-backend/internal/domain"
	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	inboundRepo repository.InboundRepository
}

func NewUserService(userRepo repository.UserRepository, inboundRepo repository.InboundRepository) *UserService {
	return &UserService{userRepo: userRepo, inboundRepo: inboundRepo}
}

func (s *UserService) List(ctx context.Context, f dto.UserFilters) (*dto.PaginatedUsers, error) {
	users, total, err := s.userRepo.List(ctx, repository.UserFilter{
		Search:   f.Search,
		Status:   f.Status,
		Protocol: f.Protocol,
		Page:     f.Page,
		PageSize: f.PageSize,
	})
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(&u)
	}

	return &dto.PaginatedUsers{
		Data:     result,
		Total:    total,
		Page:     f.Page,
		PageSize: f.PageSize,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	resp := toUserResponse(u)
	return &resp, nil
}

func (s *UserService) Create(ctx context.Context, req dto.UserFormData) (*dto.UserResponse, error) {
	trafficBytes := req.TrafficLimit * 1_000_000_000

	var expireAt *time.Time
	if req.ExpireAt != "" {
		t, err := time.Parse("2006-01-02", req.ExpireAt)
		if err == nil {
			expireAt = &t
		}
	}

	token := generateToken()

	user := &domain.User{
		Username:          req.Username,
		Email:             req.Email,
		Status:            req.Status,
		TrafficLimitBytes: trafficBytes,
		ExpireAt:          expireAt,
		MaxIPs:            req.MaxIPs,
		SubscriptionToken: token,
		Notes:             req.Notes,
	}

	if err := s.userRepo.Create(ctx, user, req.Protocols, req.InboundIDs); err != nil {
		return nil, err
	}

	resp := toUserResponse(user)
	resp.Protocols = req.Protocols
	resp.InboundIDs = req.InboundIDs
	return &resp, nil
}

func (s *UserService) Update(ctx context.Context, id string, req dto.UserFormData) (*dto.UserResponse, error) {
	existing, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	trafficBytes := req.TrafficLimit * 1_000_000_000

	var expireAt *time.Time
	if req.ExpireAt != "" {
		t, err := time.Parse("2006-01-02", req.ExpireAt)
		if err == nil {
			expireAt = &t
		}
	}

	existing.Username = req.Username
	existing.Email = req.Email
	existing.Status = req.Status
	existing.TrafficLimitBytes = trafficBytes
	existing.ExpireAt = expireAt
	existing.MaxIPs = req.MaxIPs
	existing.Notes = req.Notes

	if err := s.userRepo.Update(ctx, existing, req.Protocols, req.InboundIDs); err != nil {
		return nil, err
	}

	resp := toUserResponse(existing)
	resp.Protocols = req.Protocols
	resp.InboundIDs = req.InboundIDs
	return &resp, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ResetTraffic(ctx context.Context, id string) error {
	return s.userRepo.ResetTraffic(ctx, id)
}

func (s *UserService) ResetToken(ctx context.Context, id string) (string, error) {
	token := generateToken()
	if err := s.userRepo.UpdateSubscriptionToken(ctx, id, token); err != nil {
		return "", err
	}
	return token, nil
}

func toUserResponse(u *domain.User) dto.UserResponse {
	protocols := make([]string, len(u.Protocols))
	for i, p := range u.Protocols {
		protocols[i] = p.Protocol
	}

	return dto.UserResponse{
		ID:                u.ID,
		Username:          u.Username,
		Email:             u.Email,
		Status:            u.Status,
		Protocols:         protocols,
		InboundIDs:        u.InboundIDs,
		TrafficLimit:      u.TrafficLimitBytes,
		TrafficUsed:       u.TrafficUsedBytes,
		ExpireAt:          formatTime(u.ExpireAt),
		MaxIPs:            u.MaxIPs,
		ConcurrentIPs:     u.ConcurrentIPs,
		SubscriptionToken: u.SubscriptionToken,
		Notes:             u.Notes,
		LastSeenAt:        dto.TimePtr(nilTime(u.LastSeenAt)),
		CreatedAt:         u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         u.UpdatedAt.Format(time.RFC3339),
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func nilTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func generateToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("sub-%s", hex.EncodeToString(b)[:16])
}
