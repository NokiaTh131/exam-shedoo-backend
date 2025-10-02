package models

import "gorm.io/datatypes"

type Course struct {
	ID         int                         `gorm:"primaryKey;autoIncrement"`
	CourseCode string                      `gorm:"column:course_code;size:6;not null;uniqueIndex:idx_course"`
	Title      string                      `gorm:"column:title;not null"`
	LabSection *string                     `gorm:"column:lab_section;uniqueIndex:idx_course"`
	LecSection *string                     `gorm:"column:lec_section;uniqueIndex:idx_course"`
	Room       *string                     `gorm:"column:room"`
	LecCredit  *float32                    `gorm:"column:lec_credit"`
	LabCredit  *float32                    `gorm:"column:lab_credit"`
	Days       *string                     `gorm:"column:days"`
	StartTime  *string                     `gorm:"column:start_time"`
	EndTime    *string                     `gorm:"column:end_time"`
	Lecturers  datatypes.JSONSlice[string] `gorm:"type:json"`

	Enrollments []Enrollment `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Exams       []CourseExam `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
