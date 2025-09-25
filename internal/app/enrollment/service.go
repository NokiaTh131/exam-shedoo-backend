package enrollment

import (
	"fmt"

	"shedoo-backend/internal/dto"
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
	for i := range enrollments {
		var course models.Course
		err := s.repo.DB.Where("course_code = ? AND lec_section = ? AND lab_section = ?",
			enrollments[i].CourseCode,
			enrollments[i].LecSection,
			enrollments[i].LabSection,
		).First(&course).Error
		if err != nil {
			return fmt.Errorf("course not found for code=%s lec=%s lab=%s: %w",
				enrollments[i].CourseCode, enrollments[i].LecSection, enrollments[i].LabSection, err)
		}
		enrollments[i].CourseID = course.ID
	}
	return s.repo.BulkInsert(enrollments)
}

func (s *EnrollmentService) DeleteEnrolledByID(id uint) error {
	return s.repo.DeleteByID(id)
}

func (s *EnrollmentService) GetEnrollmentsByStudent(studentCode string) ([]dto.EnrollmentResponse, error) {
	return s.repo.GetEnrollmentsByStudent(studentCode)
}
