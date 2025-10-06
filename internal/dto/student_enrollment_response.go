package dto

type LecturerResponse struct {
	Name string `json:"name"`
}

type EnrollmentResponse struct {
	ID          int                `json:"id"`
	CourseCode  string             `json:"course_code"`
	LecSection  string             `json:"lec_section"`
	CourseName  string             `json:"course_name"`
	LabSection  string             `json:"lab_section"`
	LecCredit   float32            `json:"lec_credit"`
	LabCredit   float32            `json:"lab_credit"`
	Instructors []LecturerResponse `json:"instructors"`
	Room        string             `json:"room"`
	Days        string             `json:"days"`
	StartTime   string             `json:"start_time"`
	Semester    string             `json:"semester"`
	Year        string             `json:"year"`
	EndTime     string             `json:"end_time"`
}
