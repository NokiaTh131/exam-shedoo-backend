package models

type CourseExam struct {
	ID                   int     `gorm:"primaryKey;autoIncrement"`
	CourseCode           string  `gorm:"size:6;not null;uniqueIndex:idx_course_exam"`
	LecSection           string  `gorm:"size:10;not null;default:'000';uniqueIndex:idx_course_exam"`
	CourseID             int     `gorm:"not null"`
	Course               Course  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	LabSection           string  `gorm:"size:10;not null;default:'000';uniqueIndex:idx_course_exam"`
	MidtermExamDate      *string `gorm:"column:midterm_exam_date"`
	FinalExamDate        *string `gorm:"column:final_exam_date"`
	MidtermExamStartTime *string `gorm:"column:midterm_exam_start_time"`
	FinalExamStartTime   *string `gorm:"column:final_exam_start_time"`
	MidtermExamEndTime   *string `gorm:"column:midterm_exam_end_time"`
	FinalExamEndTime     *string `gorm:"column:final_exam_end_time"`
}
