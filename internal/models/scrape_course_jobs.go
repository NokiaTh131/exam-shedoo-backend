package models

import "time"

type ScrapeCourseJob struct {
	ID        int       `gorm:"primaryKey"`
	StartCode string    `gorm:"not null"`
	EndCode   string    `gorm:"not null"`
	Workers   int       `gorm:"default:4"`
	Status    string    `gorm:"type:varchar(20);default:'pending'"` // pending, running, completed, failed
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
