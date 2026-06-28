package xray

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/aurora/aurora-backend/internal/domain"
)

// Client — gRPC клиент для взаимодействия с Xray-core API.
// Использует HandlerService для управления inbound'ами и StatsService для сбора статистики.
type Client struct {
	addr   string
	conn   *grpc.ClientConn
	logger *zap.Logger
}

func NewClient(nodeHost string, apiPort int, logger *zap.Logger) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", nodeHost, apiPort)

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("dial xray grpc: %w", err)
	}

	return &Client{addr: addr, conn: conn, logger: logger}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Addr() string {
	return c.addr
}

// AddUser добавляет пользователя в указанный inbound.
// Xray-core поддерживает hot-reload: изменения применяются через AlterInbound без разрыва соединений.
func (c *Client) AddUser(ctx context.Context, inboundTag string, user *domain.User, protocols []string) error {
	// В реальной имплементации:
	// 1. Получаем текущий inbound через HandlerServiceClient
	// 2. Добавляем пользователя в settings.clients (VLESS/VMess/Trojan) или settings.accounts (Shadowsocks)
	// 3. Вызываем AlterInbound с обновлённым конфигом
	// Xray-core применяет hot-reload

	c.logger.Info("xray: add user stub",
		zap.String("tag", inboundTag),
		zap.String("email", user.Email),
		zap.Strings("protocols", protocols),
	)
	// TODO: implement when Xray-core proto definitions are available
	return nil
}

// RemoveUser удаляет пользователя из inbound'а.
func (c *Client) RemoveUser(ctx context.Context, inboundTag, email string) error {
	c.logger.Info("xray: remove user stub",
		zap.String("tag", inboundTag),
		zap.String("email", email),
	)
	return nil
}

// GetUserStats возвращает статистику трафика пользователя (upload/download).
func (c *Client) GetUserStats(ctx context.Context, email string) (*domain.TrafficStats, error) {
	// В реальной имплементации:
	// Используем StatsService.QueryStats с паттерном: "user>>>{email}>>>traffic>>>uplink"
	// и "user>>>{email}>>>traffic>>>downlink"

	c.logger.Debug("xray: get user stats stub", zap.String("email", email))
	return &domain.TrafficStats{Upload: 0, Download: 0}, nil
}

// GetInboundStats возвращает агрегированную статистику inbound'а.
func (c *Client) GetInboundStats(ctx context.Context, tag string) (*domain.TrafficStats, error) {
	c.logger.Debug("xray: get inbound stats stub", zap.String("tag", tag))
	return &domain.TrafficStats{Upload: 0, Download: 0}, nil
}

// QueryAllStats возвращает статистику по всем пользователям и inbound'ам.
func (c *Client) QueryAllStats(ctx context.Context) (map[string]*domain.TrafficStats, error) {
	c.logger.Debug("xray: query all stats stub")
	return make(map[string]*domain.TrafficStats), nil
}

// HealthCheck проверяет доступность Xray-core gRPC API.
func (c *Client) HealthCheck(ctx context.Context) error {
	// В реальной имплементации: пингуем StatsService или HandlerService
	return nil
}
