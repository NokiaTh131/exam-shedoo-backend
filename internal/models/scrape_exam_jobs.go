package models

import "time"

type ScrapeExamJob struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Term      string    `gorm:"size:3;not null"`
	Status    string    `gorm:"size:20;not null;default:'pending'"` // pending, running, completed, failed
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
