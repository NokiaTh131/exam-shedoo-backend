package repositories

import (
	"fmt"

	"shedoo-backend/internal/dto"
	"shedoo-backend/internal/models"

	"gorm.io/gorm"
)

type CourseExamRepository struct {
	db *gorm.DB
}

func NewCourseExamRepository(db *gorm.DB) *CourseExamRepository {
	return &CourseExamRepository{db: db}
}

func (r *CourseExamRepository) Create(exam *models.CourseExam) error {
	return r.db.Create(exam).Error
}

func (r *CourseExamRepository) GetByCourseSections(courseCode, lecSection, labSection string) (*models.CourseExam, error) {
	var exam models.CourseExam
	err := r.db.Where("course_code = ? AND lec_section = ? AND lab_section = ?", courseCode, lecSection, labSection).
		First(&exam).Error
	if err != nil {
		return nil, err
	}
	return &exam, nil
}

func (r *CourseExamRepository) FindByID(id uint) (*models.CourseExam, error) {
	var exam models.CourseExam
	if err := r.db.First(&exam, id).Error; err != nil {
		return nil, err
	}
	return &exam, nil
}

func (r *CourseExamRepository) UpdateByID(id uint, updates *models.CourseExam) error {
	respon := r.db.Model(&updates).Where("id = ?", id).Updates(updates)
	if respon.RowsAffected == 0 {
		return gorm.ErrInvalidField
	}
	return nil
}

func (r *CourseExamRepository) GetExamsByStudent(studentCode string) ([]dto.StudentExamResponse, error) {
	var exams []dto.StudentExamResponse

	err := r.db.Raw(`
        SELECT 
						e.id,
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

func (r *CourseExamRepository) GetMidtermExamReport(lecturerName string) ([]dto.LecturerCourseMidterm, error) {
	var exams []dto.LecturerCourseMidterm
	err := r.db.Raw(`
    SELECT 
        c.id AS course_id,
        c.course_code,
        c.title AS course_name,
        COALESCE(c.lec_section, '000') AS lec_section,
        COALESCE(c.lab_section, '000') AS lab_section,
        COALESCE(ce.midterm_exam_date, '') AS exam_date,
        COALESCE(ce.midterm_exam_start_time, '') AS start_time,
        COALESCE(ce.midterm_exam_end_time, '') AS end_time,
        COUNT(e.id) AS num_of_students
    FROM courses c
    LEFT JOIN LATERAL (
        SELECT *
        FROM course_exams ce
        WHERE ce.course_id = c.id
          AND (ce.lec_section = c.lec_section OR ce.lec_section = '000')
          AND (ce.lab_section = c.lab_section OR ce.lab_section = '000')
        ORDER BY (ce.lec_section = c.lec_section) DESC
        LIMIT 1
    ) ce ON true
    LEFT JOIN enrollments e
        ON e.course_id = c.id
    WHERE c.lecturers @> ?
    GROUP BY c.id, c.course_code, c.title, c.lec_section, c.lab_section,
             ce.midterm_exam_date, ce.midterm_exam_start_time, ce.midterm_exam_end_time
    ORDER BY c.course_code, c.lec_section
`, fmt.Sprintf(`["%s"]`, lecturerName)).Scan(&exams).Error

	return exams, err
}
