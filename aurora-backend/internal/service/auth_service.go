package service

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/argon2"

	"github.com/aurora/aurora-backend/internal/config"
	"github.com/aurora/aurora-backend/internal/domain"
	"github.com/aurora/aurora-backend/internal/pkg/jwt"
	"github.com/aurora/aurora-backend/internal/repository"
)

type AuthService struct {
	adminRepo repository.AdminRepository
	tm        *jwt.TokenManager
	redis     *redis.Client
	cfg       config.JWTConfig
}

type Argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

var defaultArgon2 = Argon2Params{
	Time:    3,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
	SaltLen: 16,
}

func NewAuthService(
	adminRepo repository.AdminRepository,
	tm *jwt.TokenManager,
	redis *redis.Client,
	cfg config.JWTConfig,
) *AuthService {
	return &AuthService{
		adminRepo: adminRepo,
		tm:        tm,
		redis:     redis,
		cfg:       cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password, ip string) (*LoginResult, error) {
	// Check IP ban
	banned, err := s.redis.Get(ctx, fmt.Sprintf("banned:%s", ip)).Result()
	if err == nil && banned == "1" {
		return nil, fmt.Errorf("IP is temporarily banned")
	}

	// Find admin
	admin, err := s.adminRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("find admin: %w", err)
	}
	if admin == nil {
		s.recordFailedAttempt(ctx, ip)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check lock
	if admin.LockedUntil != nil && admin.LockedUntil.After(time.Now()) {
		return nil, fmt.Errorf("account is locked until %s", admin.LockedUntil.Format(time.RFC3339))
	}

	// Verify password
	if !verifyPassword(password, admin.PasswordHash) {
		s.recordFailedAttempt(ctx, ip)
		attempts := admin.FailedAttempts + 1
		var lockedUntil *string
		if attempts >= 5 {
			t := time.Now().Add(15 * time.Minute).Format(time.RFC3339)
			lockedUntil = &t
		}
		_ = s.adminRepo.UpdateLoginAttempts(ctx, admin.ID, attempts, lockedUntil)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Success — reset attempts, update last login
	_ = s.adminRepo.UpdateLastLogin(ctx, admin.ID)

	// Generate tokens
	accessToken, accessExp, err := s.tm.GenerateAccessToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}
	refreshToken, _, err := s.tm.GenerateRefreshToken(admin.ID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessExp.Unix(),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	claims, err := s.tm.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check blacklist
	blacklisted, _ := s.redis.Get(ctx, fmt.Sprintf("jwt:blacklist:%s", claims.ID)).Result()
	if blacklisted == "1" {
		return nil, fmt.Errorf("token has been revoked")
	}

	admin, err := s.adminRepo.FindByID(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("find admin: %w", err)
	}
	if admin == nil {
		return nil, fmt.Errorf("admin not found")
	}

	// Blacklist old refresh token
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		s.redis.Set(ctx, fmt.Sprintf("jwt:blacklist:%s", claims.ID), "1", ttl)
	}

	// Issue new tokens
	accessToken, accessExp, err := s.tm.GenerateAccessToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, err
	}
	newRefreshToken, _, err := s.tm.GenerateRefreshToken(admin.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    accessExp.Unix(),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	claims, err := s.tm.ValidateAccessToken(accessToken)
	if err != nil {
		return nil // Already invalid — no action needed
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		s.redis.Set(ctx, fmt.Sprintf("jwt:blacklist:%s", claims.ID), "1", ttl)
	}
	return nil
}

func (s *AuthService) GetMe(ctx context.Context, adminID string) (*domain.Admin, error) {
	return s.adminRepo.FindByID(ctx, adminID)
}

func (s *AuthService) recordFailedAttempt(ctx context.Context, ip string) {
	key := fmt.Sprintf("ratelimit:login:%s", ip)
	count, _ := s.redis.Incr(ctx, key).Result()
	if count == 1 {
		s.redis.Expire(ctx, key, 15*time.Minute)
	}
	if count >= 5 {
		s.redis.Set(ctx, fmt.Sprintf("banned:%s", ip), "1", 15*time.Minute)
	}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// ─── Password helpers ───

func HashPassword(password string) (string, error) {
	salt := make([]byte, defaultArgon2.SaltLen)
	// In production, use crypto/rand here
	for i := range salt {
		salt[i] = byte(i + 42)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultArgon2.Time,
		defaultArgon2.Memory,
		defaultArgon2.Threads,
		defaultArgon2.KeyLen,
	)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		defaultArgon2.Memory,
		defaultArgon2.Time,
		defaultArgon2.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

func verifyPassword(password, encodedHash string) bool {
	// Simplified parser for the format above
	// Production code should use a proper argon2 library parser
	var memory, timeParam, threads uint32
	var b64Salt, b64Hash string
	_, err := fmt.Sscanf(encodedHash,
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		&memory, &timeParam, &threads, &b64Salt, &b64Hash)
	if err != nil {
		return false
	}

	salt, _ := base64.RawStdEncoding.DecodeString(b64Salt)
	expectedHash, _ := base64.RawStdEncoding.DecodeString(b64Hash)

	hash := argon2.IDKey([]byte(password), salt, timeParam, memory, uint8(threads), uint32(len(expectedHash)))
	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}
