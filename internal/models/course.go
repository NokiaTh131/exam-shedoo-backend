package models

type Course struct {
	ID         int      `gorm:"primaryKey;autoIncrement"`
	CourseCode string   `gorm:"column:course_code;size:6;not null;uniqueIndex:idx_course"`
	Title      string   `gorm:"column:title;not null"`
	LabSection *string  `gorm:"column:lab_section;uniqueIndex:idx_course"`
	LecSection *string  `gorm:"column:lec_section;uniqueIndex:idx_course"`
	Room       *string  `gorm:"column:room"`
	Credit     *float32 `gorm:"column:credit"`
	Days       *string  `gorm:"column:days"`
	StartTime  *string  `gorm:"column:start_time"`
	EndTime    *string  `gorm:"column:end_time"`
}
