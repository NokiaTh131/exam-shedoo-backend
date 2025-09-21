package repositories

import (
	"encoding/json"
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
	lecturerJSON, err := json.Marshal([]string{lecturerName})
	if err != nil {
		return nil, err
	}
	err = r.db.Where("lecturers @> ?", lecturerJSON).Find(&courses).Error
	return courses, err
}
