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
	return c.JSON(courses)
}

func (h *CourseHandler) GetEnrolledStudents(c *fiber.Ctx) error {
	courseID, err := c.ParamsInt("course_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid course ID",
		})
	}

	students, err := h.courseService.GetEnrolledStudents(uint(courseID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not get enrolled students",
		})
	}

	if len(students) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "No students enrolled in this course",
		})
	}

	return c.JSON(students)
}
