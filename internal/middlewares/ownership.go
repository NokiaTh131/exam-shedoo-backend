package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func StudentOwnsResource() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing user claims"})
		}

		studentID, _ := claims["student_id"].(string)
		requestedID := c.Params("studentCode")

		if studentID != requestedID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "you cannot access another student's data",
			})
		}

		return c.Next()
	}
}

func ProfessorOwnsResource() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing user claims"})
		}

		role, _ := c.Locals("role").(string)
		accountName, _ := claims["cmuitaccount_name"].(string)

		if role == "admin" {
			return c.Next()
		}

		if role == "professor" {
			if c.Query("lecturer") != accountName {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "you cannot access other professor's courses",
				})
			}
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}
}
