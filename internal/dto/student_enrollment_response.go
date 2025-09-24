package dto

type LecturerResponse struct {
	Name string `json:"name"`
}

type EnrollmentResponse struct {
	CourseCode  string             `json:"course_code"`
	LecSection  string             `json:"lec_section"`
	CourseName  string             `json:"course_name"`
	LabSection  string             `json:"lab_section"`
	Credit      float32            `json:"credit"`
	Instructors []LecturerResponse `json:"instructors"`
	Room        string             `json:"room"`
	Days        string             `json:"days"`
	StartTime   string             `json:"start_time"`
	EndTime     string             `json:"end_time"`
}
