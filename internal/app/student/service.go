package student

import (
	"shedoo-backend/internal/dto"
	"shedoo-backend/internal/repositories"
)

type StudentService struct {
	repo *repositories.StudentRepository
}

func NewStudentService(repo *repositories.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetEnrollmentsByStudent(studentCode string) ([]dto.EnrollmentResponse, error) {
	return s.repo.GetEnrollmentsByStudent(studentCode)
}

func (s *StudentService) GetExamsByStudent(studentCode string) ([]dto.StudentExamResponse, error) {
	return s.repo.GetExamsByStudent(studentCode)
}
