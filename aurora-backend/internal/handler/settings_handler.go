package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/repository"
)

type SettingsHandler struct {
	settingsRepo repository.SettingsRepository
}

func NewSettingsHandler(settingsRepo repository.SettingsRepository) *SettingsHandler {
	return &SettingsHandler{settingsRepo: settingsRepo}
}

func (h *SettingsHandler) GetAll(c *fiber.Ctx) error {
	settings, err := h.settingsRepo.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}

	result := make([]dto.SettingEntry, len(settings))
	for i, s := range settings {
		result[i] = dto.SettingEntry{
			Key:         s.Key,
			Value:       s.Value,
			Description: s.Description,
			UpdatedAt:   s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return c.JSON(result)
}

func (h *SettingsHandler) BulkUpdate(c *fiber.Ctx) error {
	var req dto.SettingsUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	if err := h.settingsRepo.BulkSet(c.Context(), req.Settings); err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "UPDATE_FAILED", Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Settings updated"})
}
