package scrapejobs

import (
	"shedoo-backend/internal/models"
	"shedoo-backend/internal/repositories"
)

type ScrapeJobService struct {
	repo *repositories.ScrapeJobRepository
}

func NewScrapeJobService(repo *repositories.ScrapeJobRepository) *ScrapeJobService {
	return &ScrapeJobService{repo: repo}
}

func (s *ScrapeJobService) GetCourseScrapeJobByID(id uint) (*models.ScrapeCourseJob, error) {
	return s.repo.GetCourseJobByID(id)
}

func (s *ScrapeJobService) CreateCourseScrapeJob(job *models.ScrapeCourseJob) error {
	return s.repo.CreateCourseJob(job)
}

func (s *ScrapeJobService) GetExamScrapeJobByID(id uint) (*models.ScrapeExamJob, error) {
	return s.repo.GetExamJobByID(id)
}

func (s *ScrapeJobService) CreateExamScrapeJob(job *models.ScrapeExamJob) error {
	return s.repo.CreateExamJob(job)
}
