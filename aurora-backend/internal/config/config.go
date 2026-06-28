package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Xray     XrayConfig
	Log      LogConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type XrayConfig struct {
	GRPCTimeout time.Duration
	RetryMax    int
	RetryWait   time.Duration
}

type LogConfig struct {
	Level  string
	Format string // json | console
}

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.readTimeout", "15s")
	v.SetDefault("server.writeTimeout", "30s")
	v.SetDefault("server.idleTimeout", "60s")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "aurora")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "aurora")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.maxConns", 25)
	v.SetDefault("database.minConns", 5)

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.poolSize", 20)

	v.SetDefault("jwt.accessSecret", "change-me-access-secret")
	v.SetDefault("jwt.refreshSecret", "change-me-refresh-secret")
	v.SetDefault("jwt.accessTTL", "15m")
	v.SetDefault("jwt.refreshTTL", "168h") // 7 days

	v.SetDefault("xray.grpcTimeout", "10s")
	v.SetDefault("xray.retryMax", 3)
	v.SetDefault("xray.retryWait", "2s")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	// Env overrides
	v.SetEnvPrefix("AURORA")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// Config file is optional in development
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
