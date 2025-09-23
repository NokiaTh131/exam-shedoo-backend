package models

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Firstname string `gorm:"size:100;not null"`
	Lastname  string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;not null;unique"`
	Password  string `gorm:"size:100"`
	StudentID string
	Role      string `gorm:"size:50;not null"`
}
