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

func (r *ScrapeJobRepository) CreateCourseJob(job *models.ScrapeCourseJob) error {
	return r.db.Create(&job).Error
}

func (r *ScrapeJobRepository) GetCourseJobByID(id uint) (*models.ScrapeCourseJob, error) {
	var job models.ScrapeCourseJob
	err := r.db.First(&job, id).Error
	return &job, err
}

func (r *ScrapeJobRepository) CreateExamJob(job *models.ScrapeExamJob) error {
	return r.db.Create(&job).Error
}

func (r *ScrapeJobRepository) GetExamJobByID(id uint) (*models.ScrapeExamJob, error) {
	var job models.ScrapeExamJob
	err := r.db.First(&job, id).Error
	return &job, err
}
