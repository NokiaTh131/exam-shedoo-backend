package repositories

import (
	"gorm.io/gorm"
	"shedoo-backend/internal/models"
)

type ScrapeJobRepository struct {
	db *gorm.DB
}

func NewScrapeJobRepository(db *gorm.DB) *ScrapeJobRepository {
	return &ScrapeJobRepository{db: db}
}

func (r *ScrapeJobRepository) CreateJob(job *models.ScrapeJob) error {
	return r.db.Create(&job).Error
}

func (r *ScrapeJobRepository) GetJobByID(id uint) (*models.ScrapeJob, error) {
	var job models.ScrapeJob
	err := r.db.First(&job, id).Error
	return &job, err
}
