package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

type Claims struct {
	AdminID  string `json:"admin_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwtlib.RegisteredClaims
}

func (tm *TokenManager) GenerateAccessToken(adminID, username, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(tm.accessTTL)
	claims := &Claims{
		AdminID:  adminID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			ID:        generateJTI(),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	signed, err := token.SignedString(tm.accessSecret)
	return signed, expiresAt, err
}

func (tm *TokenManager) GenerateRefreshToken(adminID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(tm.refreshTTL)
	claims := &jwtlib.RegisteredClaims{
		Subject:   adminID,
		ExpiresAt: jwtlib.NewNumericDate(expiresAt),
		IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		ID:        generateJTI(),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	signed, err := token.SignedString(tm.refreshSecret)
	return signed, expiresAt, err
}

func (tm *TokenManager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwtlib.ParseWithClaims(tokenStr, claims, func(t *jwtlib.Token) (any, error) {
		return tm.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwtlib.ErrSignatureInvalid
	}
	return claims, nil
}

func (tm *TokenManager) ValidateRefreshToken(tokenStr string) (*jwtlib.RegisteredClaims, error) {
	claims := &jwtlib.RegisteredClaims{}
	token, err := jwtlib.ParseWithClaims(tokenStr, claims, func(t *jwtlib.Token) (any, error) {
		return tm.refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwtlib.ErrSignatureInvalid
	}
	return claims, nil
}

func generateJTI() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
