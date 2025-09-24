package middlewares

import (
	"shedoo-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("shedoo-cmu-entraid-token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "missing token"})
	}
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (any, error) {
		return []byte(config.LoadAuthConfig().JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "invalid token"})
	}
	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user", claims)
	return c.Next()
}
