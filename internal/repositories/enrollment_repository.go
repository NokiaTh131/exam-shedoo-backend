package repositories

import (
	"shedoo-backend/internal/dto"
	"shedoo-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EnrollmentRepository struct {
	DB *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{DB: db}
}

func (r *EnrollmentRepository) BulkInsert(enrollments []models.Enrollment) error {
	if len(enrollments) == 0 {
		return nil
	}

	return r.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&enrollments).Error
}

func (r *EnrollmentRepository) DeleteByID(id uint) error {
	// hard delete Unscoped
	err := r.DB.Unscoped().Delete(&models.Enrollment{}, id).Error
	return err
}

func (r *EnrollmentRepository) GetEnrollmentsByStudent(studentCode string) ([]dto.EnrollmentResponse, error) {
	var enrollments []models.Enrollment
	err := r.DB.Preload("Course").Where("student_code = ?", studentCode).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	var responses []dto.EnrollmentResponse
	for _, e := range enrollments {
		var instructors []dto.LecturerResponse
		for _, l := range e.Course.Lecturers {
			instructors = append(instructors, dto.LecturerResponse{Name: l})
		}

		resp := dto.EnrollmentResponse{
			ID:          e.ID,
			CourseCode:  e.CourseCode,
			CourseName:  e.Course.Title,
			LecSection:  e.LecSection,
			LabSection:  e.LabSection,
			Credit:      derefFloat32(e.Course.Credit),
			Instructors: instructors,
			Room:        derefString(e.Course.Room),
			Days:        derefString(e.Course.Days),
			StartTime:   derefString(e.Course.StartTime),
			EndTime:     derefString(e.Course.EndTime),
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefFloat32(f *float32) float32 {
	if f != nil {
		return *f
	}
	return 0
}
