package course

import (
	"shedoo-backend/internal/dto"
	"shedoo-backend/internal/repositories"
)

type CourseService struct {
	repo *repositories.CourseRepository
}

func NewCourseService(repo *repositories.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) GetCoursesByLecturer(lecturerName string) ([]dto.ProfessorCourseResponse, error) {
	return s.repo.GetCoursesByLecturer(lecturerName)
}

func (s *CourseService) GetEnrolledStudents(courseID uint) ([]dto.CourseEnrollmentStudentResponse, error) {
	return s.repo.GetEnrolledStudents(courseID)
}
