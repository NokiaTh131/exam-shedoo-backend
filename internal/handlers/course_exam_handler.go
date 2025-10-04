package handlers

import (
	"fmt"
	"strconv"

	"shedoo-backend/internal/app/courseexam"
	"shedoo-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type CourseExamHandler struct {
	courseexamService *courseexam.CourseExamService
}

type UpdateCourseExamRequest struct {
	MidtermExamDate      *string `json:"midtermExamDate,omitempty"`
	FinalExamDate        *string `json:"finalExamDate,omitempty"`
	MidtermExamStartTime *string `json:"midtermExamStartTime,omitempty"`
	FinalExamStartTime   *string `json:"finalExamStartTime,omitempty"`
	MidtermExamEndTime   *string `json:"midtermExamEndTime,omitempty"`
	FinalExamEndTime     *string `json:"finalExamEndTime,omitempty"`
}

func NewCourseExamHandler(courseexamService *courseexam.CourseExamService) *CourseExamHandler {
	return &CourseExamHandler{courseexamService: courseexamService}
}

// POST /course_exams
func (h *CourseExamHandler) CreateExam(c *fiber.Ctx) error {
	exam := new(models.CourseExam)
	if err := c.BodyParser(exam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	created, err := h.courseexamService.CreateExam(exam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// PUT /course_exams/:id
func (h *CourseExamHandler) UpdateExam(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	exam, err := h.courseexamService.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find exam",
		})
	}

	var req UpdateCourseExamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	exam.MidtermExamDate = req.MidtermExamDate
	exam.FinalExamDate = req.FinalExamDate
	exam.MidtermExamStartTime = req.MidtermExamStartTime
	exam.FinalExamStartTime = req.FinalExamStartTime
	exam.MidtermExamEndTime = req.MidtermExamEndTime
	exam.FinalExamEndTime = req.FinalExamEndTime

	if err := h.courseexamService.UpdateExam(uint(id), exam); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "exam not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "updated examdate"})
}

func (h *CourseExamHandler) GetExams(c *fiber.Ctx) error {
	studentCode := c.Params("studentCode")

	exams, err := h.courseexamService.GetExamsByStudent(studentCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if exams == nil {
		return c.Status(404).JSON(fiber.Map{"error": "No exams found"})
	}

	return c.JSON(exams)
}

func (h *CourseExamHandler) GetExamReport(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	courseIdInt, err := strconv.Atoi(courseId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid courseId"})
	}

	reports, err := h.courseexamService.GetExamReport(courseIdInt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reports)
}

func (h *CourseExamHandler) UploadPDF(c *fiber.Ctx) error {
	// Get uploaded file
	pdfFile, err := c.FormFile("pdf")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "PDF file is required",
		})
	}

	// Get exam_type
	examType := c.FormValue("exam_type")
	if examType != "MIDTERM" && examType != "FINAL" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "exam_type must be MIDTERM or FINAL",
		})
	}

	// Save uploaded file temporarily
	tempPath := fmt.Sprintf("./tmp/%s", pdfFile.Filename)
	if err := c.SaveFile(pdfFile, tempPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save PDF",
		})
	}

	// Call service
	if err := h.courseexamService.ParseAndInsertPDF(tempPath, examType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "PDF processed successfully",
	})
}

