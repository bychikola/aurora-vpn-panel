package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"

	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/service"
)

type SubscriptionHandler struct {
	subService *service.SubscriptionService
}

func NewSubscriptionHandler(subService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: subService}
}

// ─── Private endpoints (admin) ───

func (h *SubscriptionHandler) List(c *fiber.Ctx) error {
	subs, err := h.subService.List(c.Context())
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}
	return c.JSON(subs)
}

func (h *SubscriptionHandler) Toggle(c *fiber.Ctx) error {
	sub, err := h.subService.Toggle(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "TOGGLE_FAILED", Message: err.Error(),
		})
	}
	if sub == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "Subscription not found",
		})
	}
	return c.JSON(sub)
}

func (h *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	if err := h.subService.Delete(c.Context(), c.Params("id")); err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "DELETE_FAILED", Message: err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Subscription deleted"})
}

// ─── Public endpoints (no auth, token-based) ───

// ServeSubscription обслуживает публичный эндпоинт /sub/:token
func (h *SubscriptionHandler) ServeSubscription(c *fiber.Ctx) error {
	token := c.Params("token")
	format := c.Params("format") // может быть пустым

	userAgent := c.Get("User-Agent")

	config, contentType, err := h.subService.GenerateSubscription(c.Context(), token, format, userAgent)
	if err != nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "SUBSCRIPTION_FAILED", Message: err.Error(),
		})
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "attachment; filename=aurora")
	return c.Send(config)
}

// ServeSubscriptionQRCode генерирует QR-код со ссылкой на подписку
func (h *SubscriptionHandler) ServeSubscriptionQRCode(c *fiber.Ctx) error {
	token := c.Params("token")
	subURL := "https://aurora.example.com/sub/" + token

	png, err := qrcode.Encode(subURL, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "QRCODE_FAILED", Message: err.Error(),
		})
	}

	c.Set("Content-Type", "image/png")
	return c.Send(png)
}
