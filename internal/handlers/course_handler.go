package handlers

import (
	"shedoo-backend/internal/app/course"

	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	courseService *course.CourseService
}

func NewCourseHandler(courseService *course.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

func (h *CourseHandler) GetCoursesByLecturer(c *fiber.Ctx) error {
	lecturerName := c.Query("lecturer")
	courses, err := h.courseService.GetCoursesByLecturer(lecturerName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not get courses",
		})
	}
	return c.JSON(fiber.Map{
		"courses": courses,
	})
}
