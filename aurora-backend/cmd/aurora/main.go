package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aurora/aurora-backend/internal/app"
	"github.com/aurora/aurora-backend/internal/config"
	"github.com/aurora/aurora-backend/internal/pkg/jwt"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// Logger
	zapLogger, err := zap.NewProduction(zap.AddStacktrace(zapcore.FatalLevel))
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer zapLogger.Sync()

	// Config
	cfg, err := config.Load(*configPath)
	if err != nil {
		zapLogger.Warn("failed to load config file, using defaults and env",
			zap.Error(err),
			zap.String("path", *configPath),
		)
		// Continue with defaults
		cfg = &config.Config{}
	}

	// JWT token manager
	tm := jwt.NewTokenManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)

	// Fiber app
	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: errorHandler,
		AppName:      "AURORA VPN Panel",
	})

	// Routes
	app.RegisterRoutes(fiberApp, tm, zapLogger)

	// Graceful shutdown
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	go func() {
		zapLogger.Info("starting server", zap.String("addr", addr))
		if err := fiberApp.Listen(addr); err != nil {
			zapLogger.Fatal("server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	zapLogger.Info("shutting down server", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := fiberApp.ShutdownWithContext(ctx); err != nil {
		zapLogger.Error("forced shutdown", zap.Error(err))
	}

	zapLogger.Info("server stopped")
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"code":    "INTERNAL_ERROR",
		"message": err.Error(),
	})
}
