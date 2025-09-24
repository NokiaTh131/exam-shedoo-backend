package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	return c.JSON(user)
}
