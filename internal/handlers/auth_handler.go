package handlers

import (
	"shedoo-backend/internal/app/auth"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	svc          *auth.AuthService
	cookieDomain string
	isProd       bool
}

func NewAuthHandler(svc *auth.AuthService, cookieDomain string, isProd bool) *AuthHandler {
	return &AuthHandler{svc: svc, cookieDomain: cookieDomain, isProd: isProd}
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var body struct {
		AuthorizationCode string `json:"authorizationCode"`
	}
	if err := c.BodyParser(&body); err != nil || body.AuthorizationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "invalid authorization code"})
	}
	token, err := h.svc.SignIn(c.Context(), body.AuthorizationCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "shedoo-cmu-entraid-token",
		Value:    token,
		Path:     "/",
		Domain:   h.cookieDomain,
		MaxAge:   7200,
		HTTPOnly: true,
		Secure:   h.isProd,
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{"ok": true})
}
