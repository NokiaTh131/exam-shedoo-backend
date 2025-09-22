package models

import "time"

type ScrapeCourseJob struct {
	ID        uint      `gorm:"primaryKey"`
	StartCode int       `gorm:"not null"`
	EndCode   int       `gorm:"not null"`
	Workers   int       `gorm:"default:4"`
	Status    string    `gorm:"type:varchar(20);default:'pending'"` // pending, running, completed, failed
	Progress  int       `gorm:"default:0"`
	Total     int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
