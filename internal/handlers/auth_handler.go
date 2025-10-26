// internal/handlers/auth.go
package handlers

import (
	"time"

	"shedoo-backend/internal/app/auth"
	"shedoo-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

	token, cookieExpiry, err := h.svc.SignIn(c.Context(), body.AuthorizationCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": err.Error()})
	}

	// Cookie MaxAge is seconds; Expires is used by browsers too.
	maxAge := int(time.Until(cookieExpiry).Seconds())
	c.Cookie(&fiber.Cookie{
		Name:     "shedoo-cmu-entraid-token",
		Value:    token,
		Path:     "/",
		Domain:   h.cookieDomain,
		MaxAge:   maxAge,
		Expires:  cookieExpiry,
		HTTPOnly: true,
		Secure:   h.isProd,
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{"ok": true})
}

func (h *AuthHandler) SignOut(c *fiber.Ctx) error {
	// attempt to revoke server session
	jwtToken := c.Cookies("shedoo-cmu-entraid-token")
	if jwtToken != "" {
		// parse token to get jti
		parsed, _ := jwt.Parse(jwtToken, func(t *jwt.Token) (any, error) {
			// avoid panic: Not shown here; you can parse using s.jwtSecret. If you need, update to call service.
			return []byte(h.svc.JwtSecret), nil
		})
		if parsed != nil {
			if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
				if jti, ok := claims["jti"].(string); ok && jti != "" {
					_ = h.svc.RevokeSession(c.Context(), jti)
				}
			}
		}
	}

	// clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "shedoo-cmu-entraid-token",
		Value:    "",
		Path:     "/",
		Domain:   h.cookieDomain,
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.isProd,
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{"ok": true})
}

func (h *AuthHandler) EntraIdUrl(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"ok": true, "url": config.LoadAuthConfig().EntraIdURL})
}
