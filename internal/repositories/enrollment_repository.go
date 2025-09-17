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
