package handlers

import (
	"fmt"
	"strconv"

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

// GET /enrollments/:studentCode
func (h *EnrollmentHandler) GetByStudentCode(c *fiber.Ctx) error {
    studentCode := c.Params("studentCode")
    if studentCode == "" {
        // return 400 Bad Request
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "studentCode is required",
        })
    }

    enrollments, err := h.enrollmentService.GetEnrolledByStudent(studentCode)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(enrollments)
}

// DELETE /enrollments/:id
func (h *EnrollmentHandler) DeleteByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid id",
        })
    }

    err = h.enrollmentService.DeleteEnrolledByID(uint(id))
    if err != nil {
        // อาจเช็คว่า err เป็น “not found” หรืออย่างอื่น
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "deleted successfully",
    })
}

// GET /enrollments/course?courseCode=XXX&lecSection=YYY&labSection=ZZZ
func (h *EnrollmentHandler) GetByCourseSections(c *fiber.Ctx) error {
    courseCode := c.Query("courseCode")
    lecSection := c.Query("lecSection")
    labSection := c.Query("labSection")

    if courseCode == "" || lecSection == "" || labSection == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "courseCode, lecSection, labSection are required",
        })
    }

    enrollments, err := h.enrollmentService.GetStudentsByCourseSections(courseCode, lecSection, labSection)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(enrollments)
}