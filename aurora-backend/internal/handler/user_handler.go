package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"

	"github.com/aurora/aurora-backend/internal/dto"
	"github.com/aurora/aurora-backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "15"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 15
	}

	filters := dto.UserFilters{
		Search:   c.Query("search"),
		Status:   c.Query("status", "all"),
		Protocol: c.Query("protocol", "all"),
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.userService.List(c.Context(), filters)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	user, err := h.userService.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}
	if user == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "User not found",
		})
	}
	return c.JSON(user)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.UserFormData
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	if len(req.Protocols) == 0 {
		req.Protocols = []string{"vless"}
	}
	if req.Status == "" {
		req.Status = "active"
	}

	user, err := h.userService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "CREATE_FAILED", Message: err.Error(),
		})
	}

	return c.Status(201).JSON(user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var req dto.UserFormData
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Code: "BAD_REQUEST", Message: "Invalid request body",
		})
	}

	user, err := h.userService.Update(c.Context(), c.Params("id"), req)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "UPDATE_FAILED", Message: err.Error(),
		})
	}
	if user == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "User not found",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	if err := h.userService.Delete(c.Context(), c.Params("id")); err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "DELETE_FAILED", Message: err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "User deleted"})
}

func (h *UserHandler) ResetTraffic(c *fiber.Ctx) error {
	if err := h.userService.ResetTraffic(c.Context(), c.Params("id")); err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "RESET_FAILED", Message: err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Traffic reset"})
}

func (h *UserHandler) ResetToken(c *fiber.Ctx) error {
	token, err := h.userService.ResetToken(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "RESET_FAILED", Message: err.Error(),
		})
	}
	return c.JSON(fiber.Map{"token": token})
}

func (h *UserHandler) QRCode(c *fiber.Ctx) error {
	user, err := h.userService.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "INTERNAL", Message: err.Error(),
		})
	}
	if user == nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Code: "NOT_FOUND", Message: "User not found",
		})
	}

	// Generate QR code with subscription URL
	subURL := "https://aurora.example.com/sub/" + user.SubscriptionToken

	png, err := qrcode.Encode(subURL, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Code: "QRCODE_FAILED", Message: err.Error(),
		})
	}

	c.Set("Content-Type", "image/png")
	return c.Send(png)
}
