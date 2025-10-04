package dto

type ExamReport struct {
	CourseID         int     `json:"course_id"`
	CourseCode       string  `json:"course_code"`
	CourseName       string  `json:"course_name"`
	StudentCode      string  `json:"student_code"`
	MidtermDate      *string `json:"midterm_date"`
	MidtermStartTime *string `json:"midterm_start_time"`
	MidtermEndTime   *string `json:"midterm_end_time"`
	FinalDate        *string `json:"final_date"`
	FinalStartTime   *string `json:"final_start_time"`
	FinalEndTime     *string `json:"final_end_time"`
	LecSection       *string `json:"lec_section"`
	LabSection       *string `json:"lab_section"`
}

type ProfessorReportCourseResponse struct {
	MidtermResponse *[]TermResponse `json:"midterm"`
	FinalResponse   *[]TermResponse `json:"final"`
}

type TermResponse struct {
	Date    string          `json:"date"`
	Start   string          `json:"start_time"`
	End     string          `json:"end_time"`
	Courses []CourseReponse `json:"courses"`
}

type CourseReponse struct {
	CourseID        int                `json:"course_id"`
	CourseCode      string             `json:"course_code"`
	CourseName      string             `json:"course_name"`
	LecSection      string             `json:"lec_section"`
	LabSection      string             `json:"lab_section"`
	StudentResponse *[]StudentResponse `json:"students"`
}

type StudentResponse struct {
	StudentCode string `json:"student_code"`
}
