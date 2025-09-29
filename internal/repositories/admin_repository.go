package repositories

import (
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

func ClearAllData(db *gorm.DB) error {
	if err := db.Exec("DELETE FROM enrollments").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM course_exams").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM courses").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM admins").Error; err != nil {
		return err
	}

	return nil
}
