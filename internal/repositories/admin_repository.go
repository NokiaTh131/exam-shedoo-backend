package repositories

import (
	"fmt"

	"shedoo-backend/internal/models"

	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) IsAdmin(account string) (bool, error) {
	var admin models.Admin
	err := r.DB.Where("account = ?", account).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}

func (r *AdminRepository) AddAdmin(account string) error {
	return r.DB.Create(&models.Admin{Account: account}).Error
}

func (r *AdminRepository) RemoveAdmin(account string) error {
	return r.DB.Where("account = ?", account).Delete(&models.Admin{}).Error
}

func (r *AdminRepository) ListAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	err := r.DB.Find(&admins).Error
	return admins, err
}

func (r *AdminRepository) DeleteAll() error {
	if err := r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Enrollment{}).Error; err != nil {
		return fmt.Errorf("failed to delete enrollments: %w", err)
	}

	if err := r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.CourseExam{}).Error; err != nil {
		return fmt.Errorf("failed to delete course exams: %w", err)
	}

	if err := r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Course{}).Error; err != nil {
		return fmt.Errorf("failed to delete courses: %w", err)
	}

	if err := r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.ScrapeCourseJob{}).Error; err != nil {
		return fmt.Errorf("failed to delete scrape course jobs: %w", err)
	}

	if err := r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.ScrapeExamJob{}).Error; err != nil {
		return fmt.Errorf("failed to delete scrape exam jobs: %w", err)
	}

	return nil
}
