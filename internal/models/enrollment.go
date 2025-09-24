package models

type Enrollment struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	StudentCode string `gorm:"size:9;not null;uniqueIndex:idx_enrollment_unique"`
	CourseID    int    `gorm:"not null"`
	Course      Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CourseCode  string `gorm:"size:6;not null;uniqueIndex:idx_enrollment_unique"`
	LecSection  string `gorm:"column:lec_section;uniqueIndex:idx_enrollment_unique"`
	LabSection  string `gorm:"column:lab_section;uniqueIndex:idx_enrollment_unique"`
	Semester    string `gorm:"size:2;not null"`
	Year        string `gorm:"size:4;not null"`
}
