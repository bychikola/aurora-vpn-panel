package dto

import "time"

// ─── Auth ───

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=2,max=64"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

type AdminMe struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// ─── Users ───

type UserFormData struct {
	Username     string   `json:"username" validate:"required,min=2,max=128"`
	Email        string   `json:"email" validate:"required,email,max=256"`
	Status       string   `json:"status" validate:"oneof=active disabled expired"`
	Protocols    []string `json:"protocols" validate:"required,min=1"`
	InboundIDs   []string `json:"inboundIds"`
	TrafficLimit int64    `json:"trafficLimit"` // GB, converted to bytes
	ExpireAt     string   `json:"expireAt"`
	MaxIPs       int      `json:"maxIPs" validate:"min=1,max=10"`
	Notes        string   `json:"notes"`
}

type UserResponse struct {
	ID                string   `json:"id"`
	Username          string   `json:"username"`
	Email             string   `json:"email"`
	Status            string   `json:"status"`
	Protocols         []string `json:"protocols"`
	InboundIDs        []string `json:"inboundIds"`
	TrafficLimit      int64    `json:"trafficLimit"`
	TrafficUsed       int64    `json:"trafficUsed"`
	ExpireAt          string   `json:"expireAt"`
	MaxIPs            int      `json:"maxIps"`
	ConcurrentIPs     int      `json:"concurrentIps"`
	SubscriptionToken string   `json:"subscriptionToken"`
	Notes             string   `json:"notes"`
	LastSeenAt        *string  `json:"lastSeenAt"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
}

type PaginatedUsers struct {
	Data     []UserResponse `json:"data"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type UserFilters struct {
	Search   string
	Status   string
	Protocol string
	Page     int
	PageSize int
}

// ─── Nodes ───

type NodeRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=128"`
	Host     string `json:"host" validate:"required,hostname|ip"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	APIPort  int    `json:"apiPort" validate:"required,min=1,max=65535"`
	APIKey   string `json:"apiKey" validate:"required"`
	Location string `json:"location"`
	Weight   int    `json:"weight" validate:"min=0,max=100"`
}

type NodeResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Host            string  `json:"host"`
	Port            int     `json:"port"`
	APIPort         int     `json:"apiPort"`
	Status          string  `json:"status"`
	Version         string  `json:"version"`
	CPUPercent      float64 `json:"cpuPercent"`
	MemoryPercent   float64 `json:"memoryPercent"`
	DiskPercent     float64 `json:"diskPercent"`
	UplinkSpeed     float64 `json:"uplinkSpeed"`
	DownlinkSpeed   float64 `json:"downlinkSpeed"`
	UplinkTotal     int64   `json:"uplinkTotal"`
	DownlinkTotal   int64   `json:"downlinkTotal"`
	UserCount       int     `json:"userCount"`
	InboundCount    int     `json:"inboundCount"`
	Location        string  `json:"location"`
	LastPing        string  `json:"lastPing"`
	CreatedAt       string  `json:"createdAt"`
}

// ─── Inbounds ───

type InboundRequest struct {
	NodeID         string `json:"nodeId" validate:"required,uuid"`
	Tag            string `json:"tag" validate:"required,min=2,max=128"`
	Protocol       string `json:"protocol" validate:"required,oneof=vless vmess trojan shadowsocks shadowsocks-2022 hysteria2 tuic-v5"`
	Port           int    `json:"port" validate:"required,min=1,max=65535"`
	Listen         string `json:"listen" validate:"required,ip"`
	Transport      string `json:"transport" validate:"oneof=tcp http ws grpc quic"`
	Security       string `json:"security" validate:"oneof=none tls reality xtls-vision"`
	Enable         bool   `json:"enable"`
	Settings       map[string]any `json:"settings"`
	StreamSettings map[string]any `json:"streamSettings"`
}

type InboundResponse struct {
	ID             string         `json:"id"`
	NodeID         string         `json:"nodeId"`
	NodeName       string         `json:"nodeName"`
	Tag            string         `json:"tag"`
	Protocol       string         `json:"protocol"`
	Port           int            `json:"port"`
	Listen         string         `json:"listen"`
	Transport      string         `json:"transport"`
	Security       string         `json:"security"`
	Enable         bool           `json:"enable"`
	UserCount      int            `json:"userCount"`
	Upload         int64          `json:"upload"`
	Download       int64          `json:"download"`
	Settings       map[string]any `json:"settings"`
	StreamSettings map[string]any `json:"streamSettings"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
}

// ─── Subscriptions ───

type SubscriptionResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userId"`
	Username      string  `json:"username"`
	Token         string  `json:"token"`
	URL           string  `json:"url"`
	Format        string  `json:"format"`
	Enabled       bool    `json:"enabled"`
	LastRequestAt *string `json:"lastRequestAt"`
	LastUserAgent string  `json:"userAgent"`
	CreatedAt     string  `json:"createdAt"`
}

// ─── Dashboard ───

type DashboardResponse struct {
	TotalUsers          int              `json:"totalUsers"`
	ActiveUsers         int              `json:"activeUsers"`
	TotalNodes          int              `json:"totalNodes"`
	OnlineNodes         int              `json:"onlineNodes"`
	TotalTrafficUp      int64            `json:"totalTrafficUp"`
	TotalTrafficDown    int64            `json:"totalTrafficDown"`
	ActiveConnections   int              `json:"activeConnections"`
	ProtocolDistribution []ProtocolCount `json:"protocolDistribution"`
	TrafficHistory      []TrafficPoint   `json:"trafficHistory"`
	UserGrowth          []GrowthPoint     `json:"userGrowth"`
}

type ProtocolCount struct {
	Protocol string `json:"protocol"`
	Count    int    `json:"count"`
}

type TrafficPoint struct {
	Timestamp string  `json:"timestamp"`
	Upload    float64 `json:"upload"`
	Download  float64 `json:"download"`
}

type GrowthPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// ─── Settings ───

type SettingsUpdate struct {
	Settings map[string]string `json:"settings" validate:"required"`
}

type SettingEntry struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	UpdatedAt   string `json:"updatedAt"`
}

// ─── Common ───

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TimePtr(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func TimePtrStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
