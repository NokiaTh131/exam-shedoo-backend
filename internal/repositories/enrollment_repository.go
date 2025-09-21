package repositories

import (
	"shedoo-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) BulkInsert(enrollments []models.Enrollment) error {
	if len(enrollments) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&enrollments).Error
}

func (r *EnrollmentRepository) GetByStudentCode(studentCode string) ([]models.Enrollment, error) {
    var enrollments []models.Enrollment
    err := r.db.
        Where("student_code = ?", studentCode).
        Find(&enrollments).
        Error
    return enrollments, err
}

func (r *EnrollmentRepository) DeleteByID(id uint) error {
    //hard delete Unscoped
    err := r.db.Unscoped().Delete(&models.Enrollment{}, id).Error
    return err
}

func (r *EnrollmentRepository) GetByCourseAndSections(courseCode, lecSection, labSection string) ([]models.Enrollment, error) {
    var enrollments []models.Enrollment
    err := r.db.
        Where("course_code = ? AND lec_section = ? AND lab_section = ?", courseCode, lecSection, labSection).
        Find(&enrollments).
        Error
    return enrollments, err
}

