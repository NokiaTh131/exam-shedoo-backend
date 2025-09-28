package repositories

import (
	"fmt"
	"os/exec"
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

func (r *CourseExamRepository) GetExamReport(courseId int) (dto.ProfessorReportCourseResponse, error) {
	var rows []dto.ExamReport
	err := r.db.Raw(`
		SELECT 
		c.id AS course_id,
		c.course_code,
		c.title AS course_name,
		COUNT(DISTINCT e.student_code) AS student_count,
		ce.midterm_exam_date AS midterm_date,
		ce.midterm_exam_start_time AS midterm_start_time,
		ce.midterm_exam_end_time AS midterm_end_time,
		ce.final_exam_date AS final_date,
		ce.final_exam_start_time AS final_start_time,
		ce.final_exam_end_time AS final_end_time,
		c.lec_section,
		c.lab_section
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		LEFT JOIN LATERAL (
		SELECT 
			ce.midterm_exam_date, 
			ce.midterm_exam_start_time,
			ce.midterm_exam_end_time,
			ce.final_exam_date,
			ce.final_exam_start_time,
			ce.final_exam_end_time
		FROM course_exams ce
		WHERE ce.course_code = c.course_code
		AND (
			(ce.lec_section = c.lec_section AND ce.lab_section = c.lab_section)
			OR (ce.lec_section = '000' AND ce.lab_section = '000')
		)
		ORDER BY 
			(ce.lec_section = c.lec_section AND ce.lab_section = c.lab_section) DESC
		LIMIT 1
		) ce ON TRUE
		WHERE e.student_code IN (
		SELECT student_code 
		FROM enrollments 
		WHERE course_id = ?
		)
		GROUP BY 
		c.id, c.course_code, c.title, 
		ce.midterm_exam_date, ce.midterm_exam_start_time, ce.midterm_exam_end_time,
		ce.final_exam_date, ce.final_exam_start_time, ce.final_exam_end_time,
		c.lec_section, c.lab_section
		ORDER BY student_count DESC;
`, courseId).Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		return dto.ProfessorReportCourseResponse{}, err
	}
	midtermGroups := map[string]*dto.TermResponse{}
	finalGroups := map[string]*dto.TermResponse{}

	for _, row := range rows {
		course := dto.CourseReponse{
			CourseID:      row.CourseID,
			CourseCode:    row.CourseCode,
			CourseName:    row.CourseName,
			LecSection:    derefString(row.LecSection),
			LabSection:    derefString(row.LabSection),
			NumOfStudents: row.StudentCount,
		}

		if derefString(row.MidtermDate) != "" {
			key := derefString(row.MidtermDate) + derefString(row.MidtermStartTime) + derefString(row.MidtermEndTime)
			if _, exists := midtermGroups[key]; !exists {
				midtermGroups[key] = &dto.TermResponse{
					Date:    derefString(row.MidtermDate),
					Start:   derefString(row.MidtermStartTime),
					End:     derefString(row.MidtermEndTime),
					Courses: []dto.CourseReponse{},
				}
			}
			midtermGroups[key].Courses = append(midtermGroups[key].Courses, course)
		}

		if derefString(row.FinalDate) != "" {
			key := derefString(row.FinalDate) + derefString(row.FinalStartTime) + derefString(row.FinalEndTime)
			if _, exists := finalGroups[key]; !exists {
				finalGroups[key] = &dto.TermResponse{
					Date:    derefString(row.FinalDate),
					Start:   derefString(row.FinalStartTime),
					End:     derefString(row.FinalEndTime),
					Courses: []dto.CourseReponse{},
				}
			}
			finalGroups[key].Courses = append(finalGroups[key].Courses, course)
		}
	}

	midtermResp := []dto.TermResponse{}
	for _, v := range midtermGroups {
		midtermResp = append(midtermResp, *v)
	}

	finalResp := []dto.TermResponse{}
	for _, v := range finalGroups {
		finalResp = append(finalResp, *v)
	}

	response := &dto.ProfessorReportCourseResponse{
		MidtermResponse: &midtermResp,
		FinalResponse:   &finalResp,
	}

	return *response, nil
}

func (r *CourseExamRepository) ParseAndInsert(pdfPath string, examType string) error {
    cmd := exec.Command("python3", "internal/script/examdate_from_pdf.py", "--pdf", pdfPath, "--exam_type", examType)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to run Python script: %v, output: %s", err, string(output))
    }
    fmt.Println(string(output))
    return nil
}