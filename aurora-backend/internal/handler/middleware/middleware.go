package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/aurora/aurora-backend/internal/pkg/jwt"
)

func Setup(app *fiber.App, tm *jwt.TokenManager, zapLogger *zap.Logger) {
	// Recovery
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Logger
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
	}))

	// Global rate limit (light)
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"code":    "TOO_MANY_REQUESTS",
				"message": "Rate limit exceeded. Try again later.",
			})
		},
	}))
}

// AuthRequired проверяет JWT access token в заголовке Authorization
func AuthRequired(tm *jwt.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(401).JSON(fiber.Map{
				"code":    "MISSING_TOKEN",
				"message": "Authorization header is required.",
			})
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return c.Status(401).JSON(fiber.Map{
				"code":    "INVALID_TOKEN_FORMAT",
				"message": "Authorization header must be: Bearer <token>",
			})
		}

		claims, err := tm.ValidateAccessToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"code":    "INVALID_TOKEN",
				"message": "Token is invalid or expired.",
			})
		}

		// Store admin info in context for downstream handlers
		c.Locals("adminID", claims.AdminID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// AdminOnly пропускает только админов (не readonly)
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"code":    "FORBIDDEN",
				"message": "Admin privileges required.",
			})
		}
		return c.Next()
	}
}
