package domain

import (
	"time"
)

type Admin struct {
	ID             string
	Username       string
	PasswordHash   string
	Role           string    // admin, readonly
	FailedAttempts int
	LockedUntil    *time.Time
	LastLoginAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type User struct {
	ID                string
	Username          string
	Email             string    // Xray identifier
	Status            string    // active, disabled, expired
	TrafficLimitBytes int64     // 0 = unlimited
	TrafficUsedBytes  int64
	ExpireAt          *time.Time // nil = never
	MaxIPs            int
	ConcurrentIPs     int
	SubscriptionToken string
	Notes             string
	Protocols         []Protocol
	InboundIDs        []string
	LastSeenAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Protocol struct {
	UserID   string
	Protocol string // vless, vmess, trojan, shadowsocks, shadowsocks-2022, hysteria2, tuic-v5
}

type Node struct {
	ID            string
	Name          string
	Host          string
	Port          int
	APIPort       int
	APIKey        string    // AES-encrypted
	Status        string    // online, offline, degraded
	Version       string
	Location      string
	Weight        int
	Enabled       bool
	LastPingAt    *time.Time
	Metrics       string    // JSON: cpu, ram, disk, network
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Inbound struct {
	ID             string
	NodeID         string
	Tag            string
	Protocol       string
	Port           int
	Listen         string
	Transport      string // tcp, ws, http, grpc, quic
	Security       string // none, tls, reality, xtls-vision
	Enable         bool
	Settings       string // JSONB
	StreamSettings string // JSONB
	Sniffing       string // JSONB
	UserCount      int
	UploadBytes    int64
	DownloadBytes  int64
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Joined fields
	NodeName string
}

type Subscription struct {
	ID            string
	UserID        string
	Username      string
	Token         string
	URL           string
	Format        string // base64, clash, sing-box
	Enabled       bool
	LastRequestAt *time.Time
	LastUserAgent string
	CreatedAt     time.Time
}

type TrafficStats struct {
	Upload   int64
	Download int64
}

type ConnectionLog struct {
	ID             string
	UserID         *string
	Email          string
	InboundID      *string
	NodeID         *string
	IPAddress      string
	UserAgent      string
	ConnectedAt    time.Time
	DisconnectedAt *time.Time
	UploadBytes    int64
	DownloadBytes  int64
}

type InboundTraffic struct {
	ID          int64
	InboundID   string
	UserEmail   string
	Upload      int64
	Download    int64
	CollectedAt time.Time
}

type Setting struct {
	Key         string
	Value       string // JSON
	Description string
	UpdatedAt   time.Time
}

type BannedIP struct {
	IPAddress  string
	Attempts   int
	BannedUntil time.Time
	Reason     string
	CreatedAt  time.Time
}

// Dashboard stats (computed)
type DashboardStats struct {
	TotalUsers          int
	ActiveUsers         int
	TotalNodes          int
	OnlineNodes         int
	TotalTrafficUp      int64
	TotalTrafficDown    int64
	ActiveConnections   int
	ProtocolDistribution []ProtocolCount
	TrafficHistory      []TrafficPoint
	UserGrowth          []GrowthPoint
}

type ProtocolCount struct {
	Protocol string
	Count    int
}

type TrafficPoint struct {
	Timestamp string
	Upload    float64
	Download  float64
}

type GrowthPoint struct {
	Date  string
	Count int
}
