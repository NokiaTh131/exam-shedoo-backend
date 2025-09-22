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

func (s *ScrapeJobService) GetScrapeJobByID(id uint) (*models.ScrapeJob, error) {
	return s.repo.GetJobByID(id)
}

func (s *ScrapeJobService) CreateScrapeJob(job *models.ScrapeJob) error {
	return s.repo.CreateJob(job)
}
