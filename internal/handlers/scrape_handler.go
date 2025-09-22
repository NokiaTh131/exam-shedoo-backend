package handlers

import (
	"strconv"

	scrapejobs "shedoo-backend/internal/app/scrape_jobs"
	"shedoo-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type ScrapeJobHandler struct {
	scrapeJobService *scrapejobs.ScrapeJobService
}

type ScrapeJobRequest struct {
	Start   int `json:"start"`
	End     int `json:"end"`
	Workers int `json:"workers"`
}

func NewScrapeJobHandler(scrapeJobService *scrapejobs.ScrapeJobService) *ScrapeJobHandler {
	return &ScrapeJobHandler{scrapeJobService: scrapeJobService}
}

func (h *ScrapeJobHandler) GetScrapeJobByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}
	job, err := h.scrapeJobService.GetCourseScrapeJobByID(uint(idInt))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not get scrape job",
		})
	}
	return c.JSON(job)
}

func (h *ScrapeJobHandler) CreateScrapeJob(c *fiber.Ctx) error {
	var req ScrapeJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	job := models.ScrapeCourseJob{
		StartCode: req.Start,
		EndCode:   req.End,
		Workers:   req.Workers,
		Status:    "pending",
		Total:     req.End - req.Start + 1,
	}
	if err := h.scrapeJobService.CreateCourseScrapeJob(&job); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not create scrape job",
		})
	}
	return c.JSON(fiber.Map{"job_id": job.ID, "status": "pending"})
}

func (h *ScrapeJobHandler) CreateExamScrapeJob(c *fiber.Ctx) error {
	term := c.Params("term")
	job := models.ScrapeExamJob{
		Term: term,
	}
	if err := h.scrapeJobService.CreateExamScrapeJob(&job); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not create exam scrape job",
		})
	}
	return c.JSON(fiber.Map{"job_id": job.ID, "status": "pending"})
}

func (h *ScrapeJobHandler) GetExamScrapeJobByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}
	job, err := h.scrapeJobService.GetExamScrapeJobByID(uint(idInt))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not get scrape job",
		})
	}
	return c.JSON(job)
}
