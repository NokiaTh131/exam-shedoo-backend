package dto

type LecturerCourseMidterm struct {
	CourseID      uint   `json:"course_id"`
	CourseCode    string `json:"course_code"`
	CourseName    string `json:"course_name"`
	LecSection    string `json:"lec_section"`
	LabSection    string `json:"lab_section"`
	NumOfStudents int64  `json:"number_of_relevant_students"`
	ExamDate      string `json:"exam_date"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
}

type MidtermExamReportResponse struct {
	Exams []LecturerCourseMidterm `json:"exams"`
}

type LecturerCourseFinal struct {
	CourseID      uint   `json:"course_id"`
	CourseCode    string `json:"course_code"`
	CourseName    string `json:"course_name"`
	LecSection    string `json:"lec_section"`
	LabSection    string `json:"lab_section"`
	NumOfStudents int64  `json:"number_of_relevant_students"`
	ExamDate      string `json:"exam_date"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
}

type FinalExamReportResponse struct {
	Exams []LecturerCourseFinal `json:"exams"`
}
