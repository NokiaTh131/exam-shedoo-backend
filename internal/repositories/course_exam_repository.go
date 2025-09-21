package repositories

import (
	"shedoo-backend/internal/models"

	"gorm.io/gorm"
)
type CourseExamRepository struct {
	db *gorm.DB
}

func NewCourseExamRepository(db *gorm.DB) *CourseExamRepository {
	return &CourseExamRepository{db: db}
}

func (r *CourseExamRepository) Create(exam *models.CourseExam) error {
    return r.db.Create(exam).Error
}


func (r *CourseExamRepository) GetByCourseSections(courseCode, lecSection, labSection string) (*models.CourseExam, error) {
    var exam models.CourseExam
    err := r.db.Where("course_code = ? AND lec_section = ? AND lab_section = ?", courseCode, lecSection, labSection).
        First(&exam).Error
    if err != nil {
        return nil, err
    }
    return &exam, nil
}

func (r *CourseExamRepository) UpdateByID(id uint, updates map[string]interface{}) (*models.CourseExam, error) {
    var exam models.CourseExam
    if err := r.db.First(&exam, id).Error; err != nil {
        return nil, err
    }

    if err := r.db.Model(&exam).Updates(updates).Error; err != nil {
        return nil, err
    }
    
    if err := r.db.First(&exam, id).Error; err != nil {
        return nil, err
    }
	
    return &exam, nil
}