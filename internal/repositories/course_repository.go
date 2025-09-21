package repositories

import (
	"fmt"

	"shedoo-backend/internal/models"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) GetCoursesByLecturer(lecturerName string) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Where("lecturers @> ?", fmt.Sprintf(`["%s"]`, lecturerName)).Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) GetCoursesByLecLabSection(courseCode, lecSection, labSection string) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.
		Where("course_code = ? AND lec_section = ? AND lab_section = ?", courseCode, lecSection, labSection).
		Find(&courses).
		Error
	return courses, err
}
