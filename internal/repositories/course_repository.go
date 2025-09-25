package repositories

import (
	"fmt"

	"shedoo-backend/internal/dto"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) GetCoursesByLecturer(lecturerName string) ([]dto.ProfessorCourseResponse, error) {
	var results []dto.ProfessorCourseResponse

	query := `
		SELECT 
	    c.id AS course_id,
			c.course_code,
			c.title AS course_name,
			COALESCE(c.lec_section, '000') AS lec_section,
			COALESCE(c.lab_section, '000') AS lab_section,
			ce.midterm_exam_date AS midterm_date,
			ce.midterm_exam_start_time AS midterm_start_time,
			ce.midterm_exam_end_time AS midterm_end_time,
			ce.final_exam_date AS final_date,
			ce.final_exam_start_time AS final_start_time,
			ce.final_exam_end_time AS final_end_time
		FROM courses c
		LEFT JOIN course_exams ce
			ON ce.course_code = c.course_code
			AND (ce.lec_section = c.lec_section OR ce.lec_section = '000')
			AND (ce.lab_section = c.lab_section OR ce.lab_section = '000')
		WHERE c.lecturers @> ?
		ORDER BY c.course_code, c.lec_section;
	`

	err := r.db.Raw(query, fmt.Sprintf(`["%s"]`, lecturerName)).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *CourseRepository) GetEnrolledStudents(courseID uint) ([]dto.CourseEnrollmentStudentResponse, error) {
	var students []dto.CourseEnrollmentStudentResponse

	query := `
        SELECT 
            e.id AS enrollment_id,
            e.student_code
        FROM enrollments e
        WHERE e.course_id = ?
        ORDER BY e.student_code;
    `

	err := r.db.Raw(query, courseID).Scan(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}
