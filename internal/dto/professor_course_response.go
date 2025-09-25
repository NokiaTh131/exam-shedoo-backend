package dto

type ProfessorCourseResponse struct {
	CourseID         uint    `json:"course_id"`
	CourseCode       string  `json:"course_code"`
	CourseName       string  `json:"course_name"`
	LecSection       string  `json:"lec_section"`
	LabSection       string  `json:"lab_section"`
	MidtermDate      *string `json:"midterm_date"`
	MidtermStartTime *string `json:"midterm_start_time"`
	MidtermEndTime   *string `json:"midterm_end_time"`
	FinalDate        *string `json:"final_date"`
	FinalStartTime   *string `json:"final_start_time"`
	FinalEndTime     *string `json:"final_end_time"`
}
