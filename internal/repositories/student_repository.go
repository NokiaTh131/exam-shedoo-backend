package repositories

import (
	"shedoo-backend/internal/dto"
	"shedoo-backend/internal/models"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetEnrollmentsByStudent(studentCode string) ([]dto.EnrollmentResponse, error) {
	var enrollments []models.Enrollment
	err := r.db.Preload("Course").Where("student_code = ?", studentCode).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	var responses []dto.EnrollmentResponse
	for _, e := range enrollments {
		var instructors []dto.LecturerResponse
		for _, l := range e.Course.Lecturers {
			instructors = append(instructors, dto.LecturerResponse{Name: l})
		}

		resp := dto.EnrollmentResponse{
			CourseCode:  e.CourseCode,
			CourseName:  e.Course.Title,
			LecSection:  e.LecSection,
			LabSection:  e.LabSection,
			Credit:      derefFloat32(e.Course.Credit),
			Instructors: instructors,
			Room:        derefString(e.Course.Room),
			Days:        derefString(e.Course.Days),
			StartTime:   derefString(e.Course.StartTime),
			EndTime:     derefString(e.Course.EndTime),
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func (r *StudentRepository) GetExamsByStudent(studentCode string) ([]dto.StudentExamResponse, error) {
	var exams []dto.StudentExamResponse

	err := r.db.Raw(`
        SELECT 
						e.course_code,
						c.title AS course_name,
						e.lec_section,
						e.lab_section,
						ce.midterm_exam_date AS midterm_date,
						ce.midterm_exam_start_time AS midterm_start_time,
						ce.midterm_exam_end_time AS midterm_end_time,
						ce.final_exam_date AS final_date,
						ce.final_exam_start_time AS final_start_time,
						ce.final_exam_end_time AS final_end_time
				FROM enrollments e
				JOIN courses c ON c.id = e.course_id
				LEFT JOIN course_exams ce
						ON ce.course_code = e.course_code
						AND (ce.lec_section = e.lec_section OR ce.lec_section = '000')
						AND (ce.lab_section = e.lab_section OR ce.lab_section = '000')
				WHERE e.student_code = ?
				ORDER BY e.course_code, e.lec_section;    
`, studentCode).Scan(&exams).Error
	if err != nil {
		return nil, err
	}

	return exams, nil
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefFloat32(f *float32) float32 {
	if f != nil {
		return *f
	}
	return 0
}
