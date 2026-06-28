package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps go-redis for the AURORA backend
type Client struct {
	rdb *redis.Client
}

func NewClient(addr, password string, db int, poolSize int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Underlying() *redis.Client {
	return c.rdb
}

// Traffic cache

func (c *Client) SetUserTraffic(ctx context.Context, email string, upload, download int64) error {
	pipe := c.rdb.Pipeline()
	pipe.Set(ctx, fmt.Sprintf("traffic:user:%s:upload", email), upload, 1*time.Hour)
	pipe.Set(ctx, fmt.Sprintf("traffic:user:%s:download", email), download, 1*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *Client) GetUserTraffic(ctx context.Context, email string) (upload, download int64, err error) {
	pipe := c.rdb.Pipeline()
	upCmd := pipe.Get(ctx, fmt.Sprintf("traffic:user:%s:upload", email))
	downCmd := pipe.Get(ctx, fmt.Sprintf("traffic:user:%s:download", email))
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return 0, 0, err
	}
	upload, _ = upCmd.Int64()
	download, _ = downCmd.Int64()
	return upload, download, nil
}

// Dashboard cache

func (c *Client) GetDashboardStats(ctx context.Context) (map[string]any, error) {
	data, err := c.rdb.Get(ctx, "dashboard:stats").Bytes()
	if err != nil {
		return nil, err
	}
	var stats map[string]any
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return stats, nil
}

func (c *Client) SetDashboardStats(ctx context.Context, stats map[string]any, ttl time.Duration) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, "dashboard:stats", data, ttl).Err()
}

// JWT blacklist

func (c *Client) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	return c.rdb.Set(ctx, fmt.Sprintf("jwt:blacklist:%s", jti), "1", ttl).Err()
}

func (c *Client) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	val, err := c.rdb.Get(ctx, fmt.Sprintf("jwt:blacklist:%s", jti)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

// Rate limiting

func (c *Client) IncrementLoginAttempts(ctx context.Context, ip string) (int64, error) {
	key := fmt.Sprintf("ratelimit:login:%s", ip)
	count, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		c.rdb.Expire(ctx, key, 15*time.Minute)
	}
	return count, nil
}

func (c *Client) BanIP(ctx context.Context, ip string, duration time.Duration) error {
	return c.rdb.Set(ctx, fmt.Sprintf("banned:%s", ip), "1", duration).Err()
}

func (c *Client) IsIPBanned(ctx context.Context, ip string) (bool, error) {
	val, err := c.rdb.Get(ctx, fmt.Sprintf("banned:%s", ip)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

// Node heartbeat

func (c *Client) SetNodeHeartbeat(ctx context.Context, nodeID string) error {
	return c.rdb.Set(ctx,
		fmt.Sprintf("node:heartbeat:%s", nodeID),
		time.Now().Unix(),
		2*time.Minute).Err()
}

func (c *Client) GetNodeHeartbeat(ctx context.Context, nodeID string) (time.Time, error) {
	ts, err := c.rdb.Get(ctx, fmt.Sprintf("node:heartbeat:%s", nodeID)).Int64()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(ts, 0), nil
}

// Distributed locks

func (c *Client) AcquireLock(ctx context.Context, lockName string, ttl time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, fmt.Sprintf("lock:%s", lockName), "1", ttl).Result()
}

func (c *Client) ReleaseLock(ctx context.Context, lockName string) error {
	return c.rdb.Del(ctx, fmt.Sprintf("lock:%s", lockName)).Err()
}
