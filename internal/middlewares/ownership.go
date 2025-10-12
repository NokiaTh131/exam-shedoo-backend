package middlewares

import (
	"shedoo-backend/internal/app/auth"

	"github.com/gofiber/fiber/v2"
)

func StudentOwnsResource() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawUser := c.Locals("user")
		if rawUser == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing user"})
		}

		profile, ok := rawUser.(auth.SessionData)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user type"})
		}

		studentID := profile.BasicInfo.StudentID
		requestedID := c.Params("studentCode")
		role, _ := c.Locals("role").(string)

		if role == "admin" {
			return c.Next()
		}

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
		rawUser := c.Locals("user")
		if rawUser == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing user"})
		}

		profile, ok := rawUser.(auth.SessionData)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user type"})
		}

		role, _ := c.Locals("role").(string)
		accountName := profile.BasicInfo.CmuitaccountName

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
