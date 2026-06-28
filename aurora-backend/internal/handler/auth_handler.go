package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	result, err := h.authService.Login(c.Context(), req.Username, req.Password, c.IP())
	if err != nil {
		return c.Status(401).JSON(dto.ErrorResponse{
			Code: "AUTH_FAILED", Message: err.Error(),
		})
	}

	return c.JSON(dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	result, err := h.authService.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(dto.ErrorResponse{
			Code: "REFRESH_FAILED", Message: err.Error(),
		})
	}

	return c.JSON(dto.RefreshResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := extractBearerToken(c)
	if token != "" {
		_ = h.authService.Logout(c.Context(), token)
	}
	return c.JSON(fiber.Map{"message": "Logged out"})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	adminID := c.Locals("adminID").(string)
	admin, err := h.authService.GetMe(c.Context(), adminID)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: "Failed to get admin info",
		})
	}
	if admin == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "Admin not found",
		})
	}

	return c.JSON(dto.AdminMe{
		ID:       admin.ID,
		Username: admin.Username,
		Role:     admin.Role,
	})
}

func extractBearerToken(c *fiber.Ctx) string {
	header := c.Get("Authorization")
	if len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}
	return ""
}
