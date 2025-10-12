package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"shedoo-backend/internal/app/auth"
	admin "shedoo-backend/internal/app/role"
	"shedoo-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	RoleService *admin.RoleService
	RedisClient *redis.Client
	JWTSecret   []byte
}

func NewAuthMiddleware(roleService *admin.RoleService, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		RoleService: roleService,
		RedisClient: redisClient,
		JWTSecret:   []byte(config.LoadAuthConfig().JWTSecret),
	}
}

func (a *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		cookie := c.Cookies("shedoo-cmu-entraid-token")
		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "missing token"})
		}

		// --- 1️⃣ Validate JWT signature ---
		token, err := jwt.Parse(cookie, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(a.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "invalid claims"})
		}

		// --- 2️⃣ Load session from Redis using jti ---
		jti, _ := claims["jti"].(string)
		if jti == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "missing jti in token"})
		}

		key := "session:" + jti
		sessionJSON, err := a.RedisClient.Get(ctx, key).Result()
		if err == redis.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "session not found"})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "redis error"})
		}

		var session auth.SessionData
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "invalid session data"})
		}

		// --- 3️⃣ Compute role ---
		role, err := a.RoleService.ClassifyRole(session.BasicInfo.ItaccounttypeID, session.BasicInfo.CmuitaccountName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "role classification failed"})
		}

		c.Locals("user", session)
		c.Locals("role", role)

		return c.Next()
	}
}

func (a *AuthMiddleware) RequireRoles(roles ...string) fiber.Handler {
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
