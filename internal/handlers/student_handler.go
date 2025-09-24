package handlers

import (
	"shedoo-backend/internal/app/student"

	"github.com/gofiber/fiber/v2"
)

type StudentHandler struct {
	studentService *student.StudentService
}

func NewStudentHandler(studentService *student.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

func (h *StudentHandler) GetEnrollmentsByStudent(c *fiber.Ctx) error {
	studentCode := c.Params("studentCode")
	enrollments, err := h.studentService.GetEnrollmentsByStudent(studentCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not get enrollments",
		})
	}
	return c.JSON(enrollments)
}

func (h *StudentHandler) GetExams(c *fiber.Ctx) error {
	studentCode := c.Params("studentCode")

	exams, err := h.studentService.GetExamsByStudent(studentCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(exams)
}
