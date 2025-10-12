package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	scrapejobs "shedoo-backend/internal/app/scrape_jobs"
	"shedoo-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type ScrapeJobHandler struct {
	scrapeJobService *scrapejobs.ScrapeJobService
}

type ScrapeJobRequest struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Workers int    `json:"workers"`
}

func NewScrapeJobHandler(scrapeJobService *scrapejobs.ScrapeJobService) *ScrapeJobHandler {
	return &ScrapeJobHandler{scrapeJobService: scrapeJobService}
}

func (h *ScrapeJobHandler) GetScrapeJobByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			job, err := h.scrapeJobService.GetCourseScrapeJobByID(uint(idInt))
			if err != nil {
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", "Job not found")
				w.Flush()
				return
			}

			data, _ := json.Marshal(job)
			fmt.Fprintf(w, "data: %s\n\n", data)
			err = w.Flush()
			if err != nil {
				return
			}

			if job.Status == "completed" || job.Status == "failed" {
				fmt.Fprintf(w, "event: done\ndata: %s\n\n", data)
				w.Flush()
				return
			}

			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}

func (h *ScrapeJobHandler) CreateScrapeJob(c *fiber.Ctx) error {
	var req ScrapeJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	startCodeInt, err := strconv.Atoi(req.Start)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid start code",
		})
	}
	endCodeInt, err := strconv.Atoi(req.End)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid end code",
		})
	}
	if startCodeInt > endCodeInt {
		return c.Status(400).JSON(fiber.Map{
			"error": "Start code must be less than end code",
		})
	}

	job := models.ScrapeCourseJob{
		StartCode: req.Start,
		EndCode:   req.End,
		Workers:   req.Workers,
		Status:    "pending",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			job, err := h.scrapeJobService.GetExamScrapeJobByID(uint(idInt))
			if err != nil {
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", "Job not found")
				w.Flush()
				return
			}

			data, _ := json.Marshal(job)

			fmt.Fprintf(w, "data: %s\n\n", data)
			err = w.Flush()
			if err != nil {
				fmt.Printf("SSE flush error: %v — closing connection.\n", err)
				return
			}

			if job.Status == "completed" || job.Status == "failed" {
				fmt.Fprintf(w, "event: done\ndata: %s\n\n", data)
				w.Flush()
				return
			}

			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}
