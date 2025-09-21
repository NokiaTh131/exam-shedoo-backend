package handlers

import (
	"strconv"

	"shedoo-backend/internal/app/courseexam"
	"shedoo-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type CourseExamHandler struct {
    courseexamService *courseexam.CourseExamService
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

// GET /course_exams/course?courseCode=XXX&lecSection=YYY&labSection=ZZZ
func (h *CourseExamHandler) GetByCourseSections(c *fiber.Ctx) error {
    courseCode := c.Query("courseCode")
    lecSection := c.Query("lecSection")
    labSection := c.Query("labSection")

    if courseCode == "" || lecSection == "" || labSection == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "courseCode, lecSection, labSection are required"})
    }

    exam, err := h.courseexamService.GetExamByCourseSections(courseCode, lecSection, labSection)
    if err != nil {
        
        if err == fiber.ErrNotFound || err.Error() == "record not found" {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "exam not found"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(exam)
}

// PUT /course_exams/:id
func (h *CourseExamHandler) UpdateExam(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
    }

    // body update
    type UpdateCourseExamRequest struct {
        MidtermExamDate      *string `json:"midterm_exam_date"`
        FinalExamDate        *string `json:"final_exam_date"`
        MidtermExamStartTime *string `json:"midterm_exam_start_time"`
        FinalExamStartTime   *string `json:"final_exam_start_time"`
        MidtermExamEndTime   *string `json:"midterm_exam_end_time"`
        FinalExamEndTime     *string `json:"final_exam_end_time"`
    }

    var req UpdateCourseExamRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    // map[string]interface{}  update  field is not nil
    updates := make(map[string]interface{})
    if req.MidtermExamDate != nil {
        updates["midterm_exam_date"] = req.MidtermExamDate
    }
    if req.FinalExamDate != nil {
        updates["final_exam_date"] = req.FinalExamDate
    }
    if req.MidtermExamStartTime != nil {
        updates["midterm_exam_start_time"] = req.MidtermExamStartTime
    }
    if req.FinalExamStartTime != nil {
        updates["final_exam_start_time"] = req.FinalExamStartTime
    }
    if req.MidtermExamEndTime != nil {
        updates["midterm_exam_end_time"] = req.MidtermExamEndTime
    }
    if req.FinalExamEndTime != nil {
        updates["final_exam_end_time"] = req.FinalExamEndTime
    }

    if len(updates) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no fields to update"})
    }

    updated, err := h.courseexamService.UpdateExam(uint(id), updates)
    if err != nil {
        
        if err.Error() == "record not found" {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "exam not found"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(updated)
}