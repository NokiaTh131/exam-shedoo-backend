package course

import (
	"shedoo-backend/internal/models"
	"shedoo-backend/internal/repositories"
)

type CourseService struct {
	repo *repositories.CourseRepository
}

func NewCourseService(repo *repositories.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) GetCoursesByLecturer(lecturerName string) ([]models.Course, error) {
	return s.repo.GetCoursesByLecturer(lecturerName)
}
