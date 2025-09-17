package models

type Enrollment struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	StudentCode string `gorm:"size:9;not null;uniqueIndex:idx_enrollment_unique"`
	CourseCode  string `gorm:"size:6;not null;uniqueIndex:idx_enrollment_unique"`
	LecSection  string `gorm:"column:lec_section;uniqueIndex:idx_enrollment_unique"`
	LabSection  string `gorm:"column:lab_section;uniqueIndex:idx_enrollment_unique"`
	Semester    string `gorm:"size:2;not null"`
	Year        string `gorm:"size:4;not null"`
}
