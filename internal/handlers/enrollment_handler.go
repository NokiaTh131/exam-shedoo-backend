package handlers

import (
	"fmt"

	"shedoo-backend/internal/app/enrollment"

	"github.com/gofiber/fiber/v2"
)

type EnrollmentHandler struct {
	enrollmentService *enrollment.EnrollmentService
}

func NewEnrollmentHandler(enrollmentService *enrollment.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{enrollmentService: enrollmentService}
}

func (h *EnrollmentHandler) UploadEnrollments(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "File upload required",
		})
	}

	filePath := fmt.Sprintf("./tmp/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not save file",
		})
	}
	enrollments, err := h.enrollmentService.LoadEnrollments(filePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not load enrollments",
		})
	}
	if err := h.enrollmentService.ImportEnrollments(enrollments); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not import enrollments",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Imported %d records", len(enrollments)),
	})
}
