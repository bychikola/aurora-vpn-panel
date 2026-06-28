package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/aurora/aurora-backend/internal/config"
	"github.com/aurora/aurora-backend/internal/handler"
	"github.com/aurora/aurora-backend/internal/handler/middleware"
	"github.com/aurora/aurora-backend/internal/pkg/crypto"
	"github.com/aurora/aurora-backend/internal/pkg/jwt"
	"github.com/aurora/aurora-backend/internal/repository"
	"github.com/aurora/aurora-backend/internal/service"
)

// Handlers holds all HTTP handler instances
type Handlers struct {
	Auth         *handler.AuthHandler
	Dashboard    *handler.DashboardHandler
	User         *handler.UserHandler
	Node         *handler.NodeHandler
	Inbound      *handler.InboundHandler
	Subscription *handler.SubscriptionHandler
	Settings     *handler.SettingsHandler
}

// InitServices creates all services and handlers with dependency injection
func InitServices(
	cfg *config.Config,
	tm *jwt.TokenManager,
	redisClient *redis.Client,
	logger *zap.Logger,
	adminRepo repository.AdminRepository,
	userRepo repository.UserRepository,
	nodeRepo repository.NodeRepository,
	inboundRepo repository.InboundRepository,
	subRepo repository.SubscriptionRepository,
	trafficRepo repository.TrafficRepository,
	settingsRepo repository.SettingsRepository,
) *Handlers {

	// Auth
	authService := service.NewAuthService(adminRepo, tm, redisClient, cfg.JWT)
	authHandler := handler.NewAuthHandler(authService)

	// Dashboard
	dashboardService := service.NewDashboardService(userRepo, nodeRepo, trafficRepo, redisClient)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	// Users
	userService := service.NewUserService(userRepo, inboundRepo)
	userHandler := handler.NewUserHandler(userService)

	// Nodes
	aesGCM, _ := crypto.NewAESGCM("0000000000000000000000000000000000000000000000000000000000000000")
	nodeService := service.NewNodeService(nodeRepo, inboundRepo, aesGCM)
	nodeHandler := handler.NewNodeHandler(nodeService)

	// Inbounds
	inboundService := service.NewInboundService(inboundRepo, nodeRepo, logger)
	inboundHandler := handler.NewInboundHandler(inboundService)

	// Subscriptions
	subscriptionService := service.NewSubscriptionService(subRepo, userRepo, inboundRepo, nodeRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	// Settings
	settingsHandler := handler.NewSettingsHandler(settingsRepo)

	return &Handlers{
		Auth:         authHandler,
		Dashboard:    dashboardHandler,
		User:         userHandler,
		Node:         nodeHandler,
		Inbound:      inboundHandler,
		Subscription: subscriptionHandler,
		Settings:     settingsHandler,
	}
}

// RegisterRoutes настраивает все маршруты API с реальными хендлерами
func RegisterRoutes(app *fiber.App, tm *jwt.TokenManager, log *zap.Logger, h *Handlers) {
	middleware.Setup(app, tm, log)

	api := app.Group("/api/v1")

	// Healthcheck (публичный)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "0.1.0",
		})
	})

	// Auth (публичный)
	auth := api.Group("/auth")
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.Refresh)
	auth.Post("/logout", h.Auth.Logout)
	auth.Get("/me", middleware.AuthRequired(tm), h.Auth.Me)

	// Защищённые маршруты (требуют JWT)
	protected := api.Group("", middleware.AuthRequired(tm))

	// Dashboard
	protected.Get("/dashboard/stats", h.Dashboard.Stats)

	// Users
	users := protected.Group("/users")
	users.Get("/", h.User.List)
	users.Post("/", h.User.Create)
	users.Get("/:id", h.User.GetByID)
	users.Put("/:id", h.User.Update)
	users.Delete("/:id", h.User.Delete)
	users.Post("/:id/reset-traffic", h.User.ResetTraffic)
	users.Post("/:id/reset-token", h.User.ResetToken)
	users.Get("/:id/qrcode", h.User.QRCode)

	// Nodes
	nodes := protected.Group("/nodes")
	nodes.Get("/", h.Node.List)
	nodes.Get("/:id", h.Node.GetByID)
	nodes.Post("/", h.Node.Create)
	nodes.Put("/:id", h.Node.Update)
	nodes.Delete("/:id", h.Node.Delete)

	// Inbounds
	inbounds := protected.Group("/inbounds")
	inbounds.Get("/", h.Inbound.List)
	inbounds.Get("/:id", h.Inbound.GetByID)
	inbounds.Post("/", h.Inbound.Create)
	inbounds.Put("/:id", h.Inbound.Update)
	inbounds.Delete("/:id", h.Inbound.Delete)
	inbounds.Post("/:id/reload", h.Inbound.Reload)

	// Subscriptions (private)
	subs := protected.Group("/subscriptions")
	subs.Get("/", h.Subscription.List)
	subs.Post("/:id/toggle", h.Subscription.Toggle)
	subs.Delete("/:id", h.Subscription.Delete)

	// Settings
	protected.Get("/settings", h.Settings.GetAll)
	protected.Put("/settings", h.Settings.BulkUpdate)
}

// RegisterPublicRoutes настраивает публичные эндпоинты подписок (без JWT)
func RegisterPublicRoutes(app *fiber.App, h *Handlers) {
	sub := app.Group("/api/v1/sub")
	sub.Get("/:token", h.Subscription.ServeSubscription)
	sub.Get("/:token/qrcode", h.Subscription.ServeSubscriptionQRCode)
}

// Helper types for ctx keys
type contextKey string

var (
	CtxAdminID  = contextKey("adminID")
	CtxUsername = contextKey("username")
	CtxRole     = contextKey("role")
)
