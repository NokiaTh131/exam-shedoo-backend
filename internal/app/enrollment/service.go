package enrollment

import (
	"shedoo-backend/internal/models"
	"shedoo-backend/internal/repositories"
)

type EnrollmentService struct {
	repo *repositories.EnrollmentRepository
}

func NewEnrollmentService(repo *repositories.EnrollmentRepository) *EnrollmentService {
	return &EnrollmentService{repo: repo}
}

func (s *EnrollmentService) LoadEnrollments(filePath string) ([]models.Enrollment, error) {
	return readXLSX(filePath)
}

func (s *EnrollmentService) ImportEnrollments(enrollments []models.Enrollment) error {
	return s.repo.BulkInsert(enrollments)
}
