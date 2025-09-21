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

func (h *CourseHandler) GetCourseByCodeSec(c *fiber.Ctx) error {
	courseCode := c.Query("courseCode")
	lecSection := c.Query("lecSection")
	labSection := c.Query("labSection")

	if courseCode == "" || lecSection == "" || labSection == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "courseCode, lecSection, labSection are required",
		})
	}
	courses, err := h.courseService.GetCourseByCodeSec(courseCode, lecSection, labSection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(courses)
}
