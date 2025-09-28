package handlers

import (
	admin "shedoo-backend/internal/app/role"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	RoleService *admin.RoleService
}

func NewAdminHandler(service *admin.RoleService) *AdminHandler {
	return &AdminHandler{RoleService: service}
}

func (h *AdminHandler) AddAdmin(c *fiber.Ctx) error {
	type req struct {
		Account string `json:"account"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	err := h.RoleService.AdminRepo.AddAdmin(body.Account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add admin"})
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (h *AdminHandler) RemoveAdmin(c *fiber.Ctx) error {
	account := c.Params("account")
	err := h.RoleService.AdminRepo.RemoveAdmin(account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove admin"})
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (h *AdminHandler) ListAdmins(c *fiber.Ctx) error {
	admins, err := h.RoleService.AdminRepo.ListAdmins()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list admins"})
	}
	return c.JSON(admins)
}
