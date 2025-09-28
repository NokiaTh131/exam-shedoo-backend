package models

import "time"

type Admin struct {
	ID        uint   `gorm:"primaryKey"`
	Account   string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
}
