package dto

type CourseEnrollmentStudentResponse struct {
	EnrollmentID uint   `json:"enrollment_id"`
	StudentCode  string `json:"student_code"`
}
