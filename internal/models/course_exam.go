package models

type CourseExam struct {
	CourseCode           string  `gorm:"primaryKey;column:course_code;size:6;check:length(course_code) = 6"`
	Section              *string `gorm:"column:section"`
	MidtermExamDate      *string `gorm:"column:midterm_exam_date"`
	FinalExamDate        *string `gorm:"column:final_exam_date"`
	MidtermExamStartTime *string `gorm:"column:midterm_exam_start_time"`
	FinalExamStartTime   *string `gorm:"column:final_exam_start_time"`
	MidtermExamEndTime   *string `gorm:"column:midterm_exam_end_time"`
	FinalExamEndTime     *string `gorm:"column:final_exam_end_time"`
}
