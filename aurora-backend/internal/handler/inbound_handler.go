package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/service"
)

type InboundHandler struct {
	inboundService *service.InboundService
}

func NewInboundHandler(inboundService *service.InboundService) *InboundHandler {
	return &InboundHandler{inboundService: inboundService}
}

func (h *InboundHandler) List(c *fiber.Ctx) error {
	nodeID := c.Query("nodeId")
	inbounds, err := h.inboundService.List(c.Context(), nodeID)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}
	return c.JSON(inbounds)
}

func (h *InboundHandler) GetByID(c *fiber.Ctx) error {
	inbound, err := h.inboundService.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}
	if inbound == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "Inbound not found",
		})
	}
	return c.JSON(inbound)
}

func (h *InboundHandler) Create(c *fiber.Ctx) error {
	var req dto.InboundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	inbound, err := h.inboundService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "CREATE_FAILED", Message: err.Error(),
		})
	}

	return c.Status(201).JSON(inbound)
}

func (h *InboundHandler) Update(c *fiber.Ctx) error {
	var req dto.InboundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	inbound, err := h.inboundService.Update(c.Context(), c.Params("id"), req)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "UPDATE_FAILED", Message: err.Error(),
		})
	}
	if inbound == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "Inbound not found",
		})
	}

	return c.JSON(inbound)
}

func (h *InboundHandler) Delete(c *fiber.Ctx) error {
	if err := h.inboundService.Delete(c.Context(), c.Params("id")); err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "DELETE_FAILED", Message: err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Inbound deleted"})
}

func (h *InboundHandler) Reload(c *fiber.Ctx) error {
	if err := h.inboundService.Reload(c.Context(), c.Params("id")); err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "RELOAD_FAILED", Message: err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Inbound reloaded"})
}
