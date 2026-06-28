package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/service"
)

type NodeHandler struct {
	nodeService *service.NodeService
}

func NewNodeHandler(nodeService *service.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: nodeService}
}

func (h *NodeHandler) List(c *fiber.Ctx) error {
	nodes, err := h.nodeService.List(c.Context())
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}
	return c.JSON(nodes)
}

func (h *NodeHandler) GetByID(c *fiber.Ctx) error {
	node, err := h.nodeService.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}
	if node == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "Node not found",
		})
	}
	return c.JSON(node)
}

func (h *NodeHandler) Create(c *fiber.Ctx) error {
	var req dto.NodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	node, err := h.nodeService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "CREATE_FAILED", Message: err.Error(),
		})
	}

	return c.Status(201).JSON(node)
}

func (h *NodeHandler) Update(c *fiber.Ctx) error {
	var req dto.NodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	node, err := h.nodeService.Update(c.Context(), c.Params("id"), req)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "UPDATE_FAILED", Message: err.Error(),
		})
	}
	if node == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "Node not found",
		})
	}

	return c.JSON(node)
}

func (h *NodeHandler) Delete(c *fiber.Ctx) error {
	if err := h.nodeService.Delete(c.Context(), c.Params("id")); err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "DELETE_FAILED", Message: err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Node deleted"})
}
