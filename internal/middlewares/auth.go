package middlewares

import (
	"slices"

	admin "shedoo-backend/internal/app/role"
	"shedoo-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(roleService *admin.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		accountTypeID, _ := claims["itaccounttype_id"].(string)
		accountName, _ := claims["cmuitaccount_name"].(string)

		role, err := roleService.ClassifyRole(accountTypeID, accountName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "role classification failed"})
		}
		claims["role"] = role
		c.Locals("user", claims)
		c.Locals("role", role)
		return c.Next()
	}
}

func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"ok": false, "message": "role not found",
			})
		}

		if slices.Contains(roles, role) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"ok": false, "message": "forbidden: insufficient role",
		})
	}
}
